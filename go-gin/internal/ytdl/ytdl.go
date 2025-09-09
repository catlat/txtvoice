package ytdl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Info 代表我们关心的 YouTube 基本信息
type Info struct {
	Id           string
	Title        string
	Author       string
	DurationSec  int
	Views        int64
	PublishDate  string
	ThumbnailUrl string
}

// ytDlpJSON 为 yt-dlp -J 输出中我们需要的字段
type ytDlpJSON struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	Uploader   string  `json:"uploader"`
	Channel    string  `json:"channel"`
	Duration   float64 `json:"duration"` // 支持小数秒
	ViewCount  int64   `json:"view_count"`
	UploadDate string  `json:"upload_date"`
	Thumbnails []struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"thumbnails"`
}

func getBin() string {
	if b := strings.TrimSpace(os.Getenv("YTDL_BIN")); b != "" {
		return b
	}
	return "yt-dlp"
}

func getCookies() string { return strings.TrimSpace(os.Getenv("YTDL_COOKIES_FILE")) }
func getProxy() string   { return strings.TrimSpace(os.Getenv("YTDL_PROXY")) }

func getAudioFormat() string {
	if f := strings.TrimSpace(os.Getenv("YTDL_AUDIO_FORMAT")); f != "" {
		return f
	}
	return "bestaudio[abr<=64]/worstaudio/bestaudio[abr<=96]/bestaudio"
}

// FetchInfo 使用 yt-dlp 获取视频信息（本地执行，不依赖外部服务）
func FetchInfo(ctx context.Context, idOrURL string) (*Info, error) {
	bin := getBin()
	args := []string{"-J", "--no-playlist"}
	if c := getCookies(); c != "" {
		args = append(args, "--cookies", c)
	}
	args = append(args, idOrURL)
	cctx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()
	cmd := exec.CommandContext(cctx, bin, args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp -J failed: %w", err)
	}
	var data ytDlpJSON
	if err := json.Unmarshal(out, &data); err != nil {
		return nil, fmt.Errorf("parse yt-dlp json failed: %w", err)
	}
	author := data.Uploader
	if author == "" {
		author = data.Channel
	}
	bestURL := ""
	bestArea := 0
	log.Printf("ytdl: thumbnails count=%d", len(data.Thumbnails))
	for _, t := range data.Thumbnails {
		area := t.Width * t.Height
		if area > bestArea {
			bestArea = area
			bestURL = t.URL
		}
	}
	log.Printf("ytdl: chosen thumb area=%d url=%s", bestArea, bestURL)
	return &Info{Id: data.ID, Title: data.Title, Author: author, DurationSec: int(data.Duration), Views: data.ViewCount, PublishDate: data.UploadDate, ThumbnailUrl: bestURL}, nil
}

// GetBestAudioURL 使用 yt-dlp 获取音频直链
func GetBestAudioURL(ctx context.Context, idOrURL string) (string, error) {
	bin := getBin()
	format := getAudioFormat()
	args := []string{"-f", format, "-g", "--no-playlist"}
	if c := getCookies(); c != "" {
		args = append(args, "--cookies", c)
	}
	args = append(args, idOrURL)
	log.Printf("ytdl: using audio format: %s", format)
	cctx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()
	cmd := exec.CommandContext(cctx, bin, args...)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("yt-dlp -g failed: %w", err)
	}
	url := strings.TrimSpace(string(out))
	if url == "" {
		return "", errors.New("empty audio url")
	}
	return url, nil
}

// DownloadAudioTo 直接下载音频到本地文件（不转码）
// outBase: 目标基础路径（不含扩展名），函数内部使用 .%(ext)s 模板并返回最终文件路径
func DownloadAudioTo(ctx context.Context, idOrURL, outBase string) (string, error) {
	bin := getBin()
	format := getAudioFormat()
	outTemplate := outBase + ".%(ext)s"
	args := []string{"-f", format, "--no-playlist", "-o", outTemplate}
	if c := getCookies(); c != "" {
		args = append(args, "--cookies", c)
	}
	if p := getProxy(); p != "" {
		args = append(args, "--proxy", p)
	}
	args = append(args, idOrURL)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(outBase), 0o755); err != nil {
		return "", err
	}

	log.Printf("ytdl: download to template %s format=%s", outTemplate, format)
	cctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(cctx, bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("yt-dlp download failed: %w", err)
	}
	// 查找实际文件（按 outBase.* 匹配，取最新）
	matches, _ := filepath.Glob(outBase + ".*")
	if len(matches) == 0 {
		return "", fmt.Errorf("file not found after download")
	}
	latest := matches[0]
	latestInfo, _ := os.Stat(latest)
	for _, m := range matches[1:] {
		if fi, err := os.Stat(m); err == nil {
			if latestInfo == nil || fi.ModTime().After(latestInfo.ModTime()) {
				latest = m
				latestInfo = fi
			}
		}
	}
	return latest, nil
}
