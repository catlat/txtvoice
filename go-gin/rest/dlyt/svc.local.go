package dlyt

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"go-gin/const/errcode"
	"go-gin/internal/component/db"
	"go-gin/internal/ytdl"
	"go-gin/model"
)

// LocalYtSvc 使用本地 yt-dlp 执行，避免外部 dlyt 依赖
type LocalYtSvc struct{}

func NewLocalYtSvc() IYtSvc { return &LocalYtSvc{} }

// detectVideoSource 检测视频来源
func detectVideoSource(idOrUrl string) string {
	s := strings.TrimSpace(strings.ToLower(idOrUrl))
	if s == "" {
		return "unknown"
	}

	// 哔哩哔哩链接检测
	if strings.Contains(s, "bilibili.com") || strings.Contains(s, "b23.tv") {
		return "bilibili"
	}

	// YouTube 链接检测
	if strings.Contains(s, "youtube.com") || strings.Contains(s, "youtu.be") {
		return "youtube"
	}

	// 默认当作 YouTube 处理（兼容纯 ID 输入）
	return "youtube"
}

// convertLocalImageToBase64 将本地图片文件转换为 base64 数据URL
func convertLocalImageToBase64(localPath string) (string, error) {
	// 去掉 /static/ 前缀，转换为实际文件路径
	filePath := localPath
	if strings.HasPrefix(localPath, "/static/") {
		filePath = filepath.Join("public", strings.TrimPrefix(localPath, "/static/"))
	}

	// 读取文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// 根据文件扩展名确定 MIME 类型
	ext := strings.ToLower(filepath.Ext(filePath))
	var mimeType string
	switch ext {
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".png":
		mimeType = "image/png"
	case ".webp":
		mimeType = "image/webp"
	case ".gif":
		mimeType = "image/gif"
	default:
		mimeType = "image/jpeg"
	}

	// 编码为 base64 并构造 data URL
	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}

func (s *LocalYtSvc) Info(ctx context.Context, idOrUrl string) (*InfoResp, error) {
	// 检测视频来源
	videoSource := detectVideoSource(idOrUrl)
	log.Printf("yt.info video_source=%s input=%s", videoSource, idOrUrl)

	videoId := extractVideoIDFast(idOrUrl)
	if videoId == "" {
		videoId = idOrUrl
	}

	// DB 命中直接返回；如字段缺失则回填；使用本地静态文件方案
	var video model.YoutubeVideo
	sourceSite := videoSource
	if sourceSite == "unknown" {
		sourceSite = "youtube" // 向后兼容
	}
	if err := db.WithContext(ctx).Where("source_site = ? AND video_id = ?", sourceSite, videoId).First(&video).Error(); err == nil {
		missing := strings.TrimSpace(video.ThumbnailUrl) == "" || video.PublishedAt == nil || strings.TrimSpace(video.Title) == "" || strings.TrimSpace(video.ChannelTitle) == "" || video.DurationSec == 0

		// 已有缩略图但不是本地静态链接，对于 YouTube 保存到本地并更新；哔哩哔哩保持原链接
		if !missing && strings.TrimSpace(video.ThumbnailUrl) != "" && !isLocalStaticURL(video.ThumbnailUrl) {
			if videoSource == "youtube" {
				if localURL, err2 := saveURLToLocal(ctx, video.ThumbnailUrl, buildThumbKey(video.VideoId, video.ThumbnailUrl)); err2 == nil {
					_ = db.WithContext(ctx).Model(&model.YoutubeVideo{}).Where("id = ?", video.Id).Update("thumbnail_url", localURL)
					video.ThumbnailUrl = localURL
					log.Printf("yt.info thumb_saved_local video_id=%s url=%s", video.VideoId, localURL)
				} else {
					log.Printf("yt.info thumb_save_local_failed video_id=%s err=%v", video.VideoId, err2)
				}
			} else {
				log.Printf("yt.info bilibili_thumb_keep_original video_id=%s url=%s", video.VideoId, video.ThumbnailUrl)
			}
		}

		if !missing {
			log.Printf("yt.info db_hit ok video_id=%s thumb_ok=%v", video.VideoId, strings.TrimSpace(video.ThumbnailUrl) != "")

			// 对于 YouTube，如果是本地静态文件，转换为 base64
			thumbnailUrl := video.ThumbnailUrl
			if videoSource == "youtube" && isLocalStaticURL(video.ThumbnailUrl) {
				if base64URL, err := convertLocalImageToBase64(video.ThumbnailUrl); err == nil {
					thumbnailUrl = base64URL
					log.Printf("yt.info thumb_converted_to_base64 video_id=%s", video.VideoId)
				} else {
					log.Printf("yt.info thumb_convert_base64_failed video_id=%s err=%v", video.VideoId, err)
				}
			}

			return &InfoResp{
				Id:           video.VideoId,
				Title:        video.Title,
				Author:       video.ChannelTitle,
				DurationSec:  video.DurationSec,
				Views:        0,
				PublishDate:  stringOrEmpty(video.PublishedAt),
				ThumbnailUrl: thumbnailUrl,
			}, nil
		}

		log.Printf("yt.info db_hit but_missing video_id=%s need_fill_thumb=%v need_fill_pub=%v", video.VideoId, strings.TrimSpace(video.ThumbnailUrl) == "", video.PublishedAt == nil)
		info, ferr := ytdl.FetchInfo(ctx, idOrUrl)
		if ferr != nil {
			log.Printf("yt.info backfill_fetch_failed video_id=%s err=%v", video.VideoId, ferr)

			// 对于 YouTube，如果是本地静态文件，转换为 base64
			thumbnailUrl := video.ThumbnailUrl
			if videoSource == "youtube" && isLocalStaticURL(video.ThumbnailUrl) {
				if base64URL, err := convertLocalImageToBase64(video.ThumbnailUrl); err == nil {
					thumbnailUrl = base64URL
					log.Printf("yt.info thumb_converted_to_base64 video_id=%s", video.VideoId)
				} else {
					log.Printf("yt.info thumb_convert_base64_failed video_id=%s err=%v", video.VideoId, err)
				}
			}

			return &InfoResp{
				Id:           video.VideoId,
				Title:        video.Title,
				Author:       video.ChannelTitle,
				DurationSec:  video.DurationSec,
				Views:        0,
				PublishDate:  stringOrEmpty(video.PublishedAt),
				ThumbnailUrl: thumbnailUrl,
			}, nil
		}
		thumbURL := info.ThumbnailUrl
		if strings.TrimSpace(thumbURL) != "" {
			if videoSource == "youtube" {
				if localURL, upErr := saveURLToLocal(ctx, thumbURL, buildThumbKey(info.Id, thumbURL)); upErr == nil {
					thumbURL = localURL
					log.Printf("yt.info backfill_thumb_saved video_id=%s url=%s", info.Id, thumbURL)
				} else {
					log.Printf("yt.info backfill_thumb_save_failed video_id=%s err=%v", info.Id, upErr)
				}
			} else {
				log.Printf("yt.info backfill_bilibili_thumb_keep_original video_id=%s url=%s", info.Id, thumbURL)
			}
		}
		updates := map[string]any{}
		if strings.TrimSpace(video.Title) == "" && strings.TrimSpace(info.Title) != "" {
			updates["title"] = info.Title
		}
		if strings.TrimSpace(video.ChannelTitle) == "" && strings.TrimSpace(info.Author) != "" {
			updates["channel_title"] = info.Author
		}
		if video.DurationSec == 0 && info.DurationSec > 0 {
			updates["duration_sec"] = info.DurationSec
		}
		if video.PublishedAt == nil && strings.TrimSpace(info.PublishDate) != "" {
			updates["published_at"] = info.PublishDate
		}
		if strings.TrimSpace(video.ThumbnailUrl) == "" && strings.TrimSpace(thumbURL) != "" {
			updates["thumbnail_url"] = thumbURL
		}
		if len(updates) > 0 {
			_ = db.WithContext(ctx).Model(&model.YoutubeVideo{}).Where("id = ?", video.Id).Updates(updates)
			// 同步结构体的缩略图，便于返回
			if v, ok := updates["thumbnail_url"].(string); ok {
				video.ThumbnailUrl = v
			}
			log.Printf("yt.info db_backfilled video_id=%s fields=%v", video.VideoId, updates)
		}
		// 对于 YouTube，如果是本地静态文件，转换为 base64
		finalThumbnailUrl := chooseNonEmpty(video.ThumbnailUrl, thumbURL)
		if videoSource == "youtube" && isLocalStaticURL(finalThumbnailUrl) {
			if base64URL, err := convertLocalImageToBase64(finalThumbnailUrl); err == nil {
				finalThumbnailUrl = base64URL
				log.Printf("yt.info backfill_thumb_converted_to_base64 video_id=%s", video.VideoId)
			} else {
				log.Printf("yt.info backfill_thumb_convert_base64_failed video_id=%s err=%v", video.VideoId, err)
			}
		}

		return &InfoResp{
			Id:           video.VideoId,
			Title:        coalesce(video.Title, info.Title),
			Author:       coalesce(video.ChannelTitle, info.Author),
			DurationSec:  ternInt(video.DurationSec, info.DurationSec),
			Views:        0,
			PublishDate:  coalesce(stringOrEmpty(video.PublishedAt), info.PublishDate),
			ThumbnailUrl: finalThumbnailUrl,
		}, nil
	}

	// 未命中调用 yt-dlp
	log.Printf("yt.info db_miss fetch video_id=%s", videoId)
	info, err := ytdl.FetchInfo(ctx, idOrUrl)
	if err != nil {
		log.Printf("yt.info fetch_failed input=%s err=%v", idOrUrl, err)
		return nil, errcode.ErrDLYTUpstream
	}

	// 对于 YouTube 下载缩略图到本地；哔哩哔哩保持原链接
	thumbURL := info.ThumbnailUrl
	if strings.TrimSpace(thumbURL) != "" {
		if videoSource == "youtube" {
			if localURL, upErr := saveURLToLocal(ctx, thumbURL, buildThumbKey(info.Id, thumbURL)); upErr == nil {
				thumbURL = localURL
				log.Printf("yt.info thumb_saved video_id=%s url=%s", info.Id, thumbURL)
			} else {
				log.Printf("yt.info thumb_save_failed video_id=%s err=%v", info.Id, upErr)
			}
		} else {
			log.Printf("yt.info bilibili_thumb_keep_original video_id=%s url=%s", info.Id, thumbURL)
		}
	}

	// 回写 DB
	v := model.YoutubeVideo{
		SourceSite:   sourceSite,
		VideoId:      info.Id,
		Title:        info.Title,
		ChannelTitle: info.Author,
		DurationSec:  info.DurationSec,
		PublishedAt:  toPtr(info.PublishDate),
		ThumbnailUrl: thumbURL,
	}
	_ = db.WithContext(ctx).Create(&v)
	log.Printf("yt.info db_created video_id=%s", info.Id)

	// 对于 YouTube，如果是本地静态文件，转换为 base64
	finalThumbURL := thumbURL
	if videoSource == "youtube" && isLocalStaticURL(thumbURL) {
		if base64URL, err := convertLocalImageToBase64(thumbURL); err == nil {
			finalThumbURL = base64URL
			log.Printf("yt.info thumb_converted_to_base64 video_id=%s", info.Id)
		} else {
			log.Printf("yt.info thumb_convert_base64_failed video_id=%s err=%v", info.Id, err)
		}
	}

	return &InfoResp{
		Id:           info.Id,
		Title:        info.Title,
		Author:       info.Author,
		DurationSec:  info.DurationSec,
		Views:        info.Views,
		PublishDate:  info.PublishDate,
		ThumbnailUrl: finalThumbURL,
	}, nil
}

func (s *LocalYtSvc) Audio(ctx context.Context, idOrUrl string) (*AudioResp, error) {
	videoId := extractVideoIDFast(idOrUrl)
	if videoId == "" {
		videoId = idOrUrl
	}

	// DB 命中音频直链则直接返回；如非本地静态链接则使用 yt-dlp 直接下载到本地
	var video model.YoutubeVideo
	if err := db.WithContext(ctx).Where("source_site = ? AND video_id = ?", "youtube", videoId).First(&video).Error(); err == nil {
		if strings.TrimSpace(video.AudioUrl) != "" {
			log.Printf("yt.audio db_hit video_id=%s current_url=%s is_local=%v", video.VideoId, video.AudioUrl, isLocalStaticURL(video.AudioUrl))
			if !isLocalStaticURL(video.AudioUrl) {
				outBase := filepath.Join("public", "yt", "audio", video.VideoId)
				if localFile, upErr := ytdl.DownloadAudioTo(ctx, idOrUrl, outBase); upErr == nil {
					localURL := localPathToStatic(localFile)
					_ = db.WithContext(ctx).Model(&model.YoutubeVideo{}).Where("id = ?", video.Id).Update("audio_url", localURL)
					video.AudioUrl = localURL
					log.Printf("yt.audio mirrored_to_local video_id=%s url=%s", video.VideoId, localURL)
				} else {
					log.Printf("yt.audio mirror_local_failed video_id=%s err=%v", video.VideoId, upErr)
				}
			}
			return &AudioResp{Id: video.VideoId, Title: video.Title, AudioUrl: video.AudioUrl}, nil
		}
	}

	// 未命中则直接下载到本地并返回静态路径
	outBase := filepath.Join("public", "yt", "audio", videoId)
	localFile, err := ytdl.DownloadAudioTo(ctx, idOrUrl, outBase)
	if err != nil {
		log.Printf("yt.audio download_failed input=%s err=%v", idOrUrl, err)
		return nil, errcode.ErrDLYTUpstream
	}
	finalAudioURL := localPathToStatic(localFile)
	log.Printf("yt.audio saved_local video_id=%s url=%s", videoId, finalAudioURL)

	// 写回 DB
	if video.Id != 0 {
		_ = db.WithContext(ctx).Model(&model.YoutubeVideo{}).Where("id = ?", video.Id).Updates(map[string]any{"audio_url": finalAudioURL})
		log.Printf("yt.audio db_updated video_id=%s", video.VideoId)
	} else {
		_ = db.WithContext(ctx).Create(&model.YoutubeVideo{SourceSite: "youtube", VideoId: videoId, AudioUrl: finalAudioURL})
		log.Printf("yt.audio db_created video_id=%s", videoId)
	}

	return &AudioResp{Id: videoId, Title: video.Title, AudioUrl: finalAudioURL}, nil
}

// 工具函数
func extractVideoIDFast(input string) string {
	s := strings.TrimSpace(input)
	if s == "" {
		return ""
	}
	if !strings.Contains(s, "http") {
		return s
	}
	r1 := regexp.MustCompile(`(?i)youtu\.be/([A-Za-z0-9_-]{6,})`)
	if m := r1.FindStringSubmatch(s); len(m) == 2 {
		return m[1]
	}
	r2 := regexp.MustCompile(`(?i)[?&]v=([A-Za-z0-9_-]{6,})`)
	if m := r2.FindStringSubmatch(s); len(m) == 2 {
		return m[1]
	}
	return ""
}

func toPtr(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	v := s
	return &v
}

func stringOrEmpty(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func coalesce(a string, b string) string {
	if strings.TrimSpace(a) != "" {
		return a
	}
	return b
}

func chooseNonEmpty(a string, b string) string {
	if strings.TrimSpace(a) != "" {
		return a
	}
	return b
}

func ternInt(a int, b int) int {
	if a > 0 {
		return a
	}
	return b
}

func isLocalStaticURL(u string) bool {
	u = strings.TrimSpace(u)
	return strings.HasPrefix(u, "/static/")
}

func localPathToStatic(localPath string) string {
	unix := filepath.ToSlash(localPath)
	if strings.HasPrefix(unix, "public/") {
		return "/static/" + strings.TrimPrefix(unix, "public/")
	}
	return "/static/" + unix // best effort
}

func saveURLToLocal(ctx context.Context, remoteURL string, key string) (string, error) {
	// local filesystem path ./public/{key}
	localPath := filepath.Join("public", filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(localPath), 0o755); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, remoteURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return "", errcode.ErrDLYTUpstream
	}

	out, err := os.Create(localPath)
	if err != nil {
		return "", err
	}
	defer out.Close()
	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", err
	}
	// Return public URL
	return "/static/" + strings.TrimLeft(filepath.ToSlash(key), "/"), nil
}

func buildThumbKey(vid string, src string) string {
	ext := "jpg"
	if i := strings.LastIndex(src, "."); i > 0 && i+1 < len(src) {
		ext = strings.ToLower(src[i+1:])
	}
	return "yt/thumb/" + vid + "." + ext
}

func buildAudioKey(vid string, src string) string {
	// 不可靠的扩展名，给缺省 m4a
	ext := "m4a"
	if i := strings.LastIndex(src, "."); i > 0 && i+1 < len(src) {
		val := strings.ToLower(src[i+1:])
		if len(val) <= 4 {
			ext = val
		}
	}
	return "yt/audio/" + vid + "." + ext
}
