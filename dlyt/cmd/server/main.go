package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"dlyt/internal/qiniu"
	"dlyt/internal/ytdl"
)

type youtubeRequest struct {
	IDOrURL string `json:"id_or_url"`
}

type youtubeInfoResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	DurationSec int64  `json:"duration_sec"`
	Views       int64  `json:"views"`
	PublishDate string `json:"publish_date"`
	Thumbnail   string `json:"thumbnail_url"`
}

type youtubeAudioResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	AudioURL string `json:"audio_url"`
}

func main() {
	addr := ":9005"
	if v := strings.TrimSpace(os.Getenv("PORT")); v != "" {
		addr = ":" + v
	}

	uploader := qiniu.NewUploader(qiniu.Config{
		AccessKey: "yWGCc-1s2A9qBiei5TZIcRykGaoKz9YUfRXFw2o4",
		SecretKey: "Wnurh4OquhZgATs3meDJJqZiST9-LFq_YEOV-dLm",
		Bucket:    "ytb-dl",
		Domain:    "https://oss.duckai.cn", // 不带尾部斜杠
	})

	// 基本信息接口：获取视频信息 + 缩略图
	http.HandleFunc("/api/yt/info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte("method not allowed"))
			return
		}
		var req youtubeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.IDOrURL) == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("bad request: need id_or_url"))
			return
		}

		ctx := r.Context()
		info, bestThumb, _, _, _, err := ytdl.FetchBasic(ctx, req.IDOrURL)
		if err != nil {
			log.Printf("fetch error: %v", err)
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte("upstream error"))
			return
		}

		thumbURL := ""
		if bestThumb.URL != "" {
			// 下载缩略图并上传
			ru, rsize, rmime, terr := ytdl.DownloadURL(ctx, bestThumb.URL)
			if terr == nil {
				defer ru.Close()
				key := fmt.Sprintf("thumbs/%s_%dx%d.jpg", info.ID, bestThumb.Width, bestThumb.Height)
				mime := rmime
				if mime == "" {
					mime = "image/jpeg"
				}
				upURL, uerr := uploader.Put(ctx, key, ru, rsize, mime)
				if uerr == nil {
					thumbURL = upURL
				}
			}
		}

		resp := youtubeInfoResponse{
			ID:          info.ID,
			Title:       info.Title,
			Author:      info.Author,
			DurationSec: int64(info.Duration.Seconds()),
			Views:       int64(info.Views),
			PublishDate: info.PublishDate.Format("2006-01-02"),
			Thumbnail:   thumbURL,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})

	// 音频接口：仅下载并上传音频
	http.HandleFunc("/api/yt/audio", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte("method not allowed"))
			return
		}
		var req youtubeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.IDOrURL) == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("bad request: need id_or_url"))
			return
		}

		ctx := r.Context()
		info, _, audioReader, audioName, audioSize, err := ytdl.FetchBasic(ctx, req.IDOrURL)
		if err != nil {
			log.Printf("fetch error: %v", err)
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte("upstream error"))
			return
		}

		audioURL := ""
		if audioReader != nil {
			defer audioReader.Close()
			key := fmt.Sprintf("audio/%s", audioName)
			lower := strings.ToLower(audioName)
			mime := "application/octet-stream"
			switch {
			case strings.HasSuffix(lower, ".webm"):
				mime = "audio/webm"
			case strings.HasSuffix(lower, ".m4a"):
				mime = "audio/mp4"
			case strings.HasSuffix(lower, ".mp3"):
				mime = "audio/mpeg"
			case strings.HasSuffix(lower, ".3gp"):
				mime = "audio/3gpp"
			}
			upURL, uerr := uploader.Put(ctx, key, audioReader, audioSize, mime)
			if uerr == nil {
				audioURL = upURL
			}
		}

		resp := youtubeAudioResponse{
			ID:       info.ID,
			Title:    info.Title,
			AudioURL: audioURL,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})

	s := &http.Server{Addr: addr, ReadTimeout: 15 * time.Second, WriteTimeout: 60 * time.Second}
	log.Printf("listening on %s", addr)
	log.Printf("endpoints: /api/yt/info (基本信息+缩略图), /api/yt/audio (仅音频)")
	log.Fatal(s.ListenAndServe())
}
