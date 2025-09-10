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

	// 哔哩哔哩 BV 号检测（BV开头的11位字符串）
	if regexp.MustCompile(`^bv[a-z0-9]{9}$`).MatchString(s) {
		return "bilibili"
	}

	// 哔哩哔哩 av 号检测（av开头的数字）
	if regexp.MustCompile(`^av\d+$`).MatchString(s) {
		return "bilibili"
	}

	// YouTube 链接检测
	if strings.Contains(s, "youtube.com") || strings.Contains(s, "youtu.be") {
		return "youtube"
	}

	// YouTube ID 检测（11位字符串，包含字母数字和-_）
	if regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`).MatchString(s) {
		return "youtube"
	}

	// 默认当作 YouTube 处理（向后兼容）
	return "youtube"
}

// constructFullURL 根据视频来源和ID构造完整的URL
func constructFullURL(idOrUrl, videoSource string) string {
	// 如果已经是完整的URL，直接返回
	if strings.Contains(idOrUrl, "http") {
		return idOrUrl
	}

	// 根据视频来源构造URL
	switch videoSource {
	case "bilibili":
		if strings.HasPrefix(strings.ToLower(idOrUrl), "bv") {
			return "https://www.bilibili.com/video/" + idOrUrl
		}
		if strings.HasPrefix(strings.ToLower(idOrUrl), "av") {
			return "https://www.bilibili.com/video/" + idOrUrl
		}
	case "youtube":
		// YouTube 的情况保持原样
		return idOrUrl
	}

	return idOrUrl
}

// extractVideoID 提取视频ID，支持YouTube和Bilibili
func extractVideoID(input, videoSource string) string {
	s := strings.TrimSpace(input)
	if s == "" {
		return ""
	}

	switch videoSource {
	case "bilibili":
		// Bilibili BV号或av号
		if !strings.Contains(s, "http") {
			return s // 直接返回BV号或av号
		}
		// 从URL中提取BV号
		if m := regexp.MustCompile(`(?i)/video/(BV[A-Za-z0-9]+)`).FindStringSubmatch(s); len(m) == 2 {
			return m[1]
		}
		// 从URL中提取av号
		if m := regexp.MustCompile(`(?i)/video/(av\d+)`).FindStringSubmatch(s); len(m) == 2 {
			return m[1]
		}
	case "youtube":
		// YouTube 使用原有的提取逻辑
		return extractVideoIDFast(s)
	}

	return s
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
	return s.InfoWithPlatform(ctx, idOrUrl, "")
}

func (s *LocalYtSvc) InfoWithPlatform(ctx context.Context, idOrUrl, platform string) (*InfoResp, error) {
	// 优先使用前端传递的平台类型，否则自动检测
	videoSource := platform
	if videoSource == "" {
		videoSource = detectVideoSource(idOrUrl)
	}
	log.Printf("yt.info video_source=%s (from_frontend=%t) input=%s", videoSource, platform != "", idOrUrl)

	// 使用新的extractVideoID函数
	videoId := extractVideoID(idOrUrl, videoSource)
	if videoId == "" {
		videoId = idOrUrl
	}

	// 构造完整URL供yt-dlp使用
	fullURL := constructFullURL(idOrUrl, videoSource)

	// DB 命中直接返回；如字段缺失则回填；使用本地静态文件方案
	var video model.YoutubeVideo
	sourceSite := videoSource
	if sourceSite == "unknown" {
		sourceSite = "youtube" // 向后兼容
	}
	if err := db.WithContext(ctx).Where("source_site = ? AND video_id = ?", sourceSite, videoId).First(&video).Error(); err == nil {
		missing := strings.TrimSpace(video.ThumbnailUrl) == "" || video.PublishedAt == nil || strings.TrimSpace(video.Title) == "" || strings.TrimSpace(video.ChannelTitle) == "" || video.DurationSec == 0

		// 已有缩略图但不是本地静态链接，统一保存到本地
		if !missing && strings.TrimSpace(video.ThumbnailUrl) != "" && !isLocalStaticURL(video.ThumbnailUrl) {
			if localURL, err2 := saveURLToLocal(ctx, video.ThumbnailUrl, buildThumbKey(video.VideoId, video.ThumbnailUrl)); err2 == nil {
				_ = db.WithContext(ctx).Model(&model.YoutubeVideo{}).Where("id = ?", video.Id).Update("thumbnail_url", localURL)
				video.ThumbnailUrl = localURL
				log.Printf("yt.info thumb_saved_local video_id=%s url=%s source=%s", video.VideoId, localURL, videoSource)
			} else {
				log.Printf("yt.info thumb_save_local_failed video_id=%s err=%v source=%s", video.VideoId, err2, videoSource)
			}
		}

		if !missing {
			log.Printf("yt.info db_hit ok video_id=%s thumb_ok=%v", video.VideoId, strings.TrimSpace(video.ThumbnailUrl) != "")

			// 如果是本地静态文件，统一转换为 base64
			thumbnailUrl := video.ThumbnailUrl
			if isLocalStaticURL(video.ThumbnailUrl) {
				if base64URL, err := convertLocalImageToBase64(video.ThumbnailUrl); err == nil {
					thumbnailUrl = base64URL
					log.Printf("yt.info thumb_converted_to_base64 video_id=%s source=%s", video.VideoId, videoSource)
				} else {
					log.Printf("yt.info thumb_convert_base64_failed video_id=%s err=%v source=%s", video.VideoId, err, videoSource)
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
		info, ferr := ytdl.FetchInfoWithPlatform(ctx, fullURL, videoSource)
		if ferr != nil {
			log.Printf("yt.info backfill_fetch_failed video_id=%s err=%v", video.VideoId, ferr)

			// 如果是本地静态文件，统一转换为 base64
			thumbnailUrl := video.ThumbnailUrl
			if isLocalStaticURL(video.ThumbnailUrl) {
				if base64URL, err := convertLocalImageToBase64(video.ThumbnailUrl); err == nil {
					thumbnailUrl = base64URL
					log.Printf("yt.info thumb_converted_to_base64 video_id=%s source=%s", video.VideoId, videoSource)
				} else {
					log.Printf("yt.info thumb_convert_base64_failed video_id=%s err=%v source=%s", video.VideoId, err, videoSource)
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
			if localURL, upErr := saveURLToLocal(ctx, thumbURL, buildThumbKey(info.Id, thumbURL)); upErr == nil {
				thumbURL = localURL
				log.Printf("yt.info backfill_thumb_saved video_id=%s url=%s source=%s", info.Id, thumbURL, videoSource)
			} else {
				log.Printf("yt.info backfill_thumb_save_failed video_id=%s err=%v source=%s", info.Id, upErr, videoSource)
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
		// 如果是本地静态文件，统一转换为 base64
		finalThumbnailUrl := chooseNonEmpty(video.ThumbnailUrl, thumbURL)
		if isLocalStaticURL(finalThumbnailUrl) {
			if base64URL, err := convertLocalImageToBase64(finalThumbnailUrl); err == nil {
				finalThumbnailUrl = base64URL
				log.Printf("yt.info backfill_thumb_converted_to_base64 video_id=%s source=%s", video.VideoId, videoSource)
			} else {
				log.Printf("yt.info backfill_thumb_convert_base64_failed video_id=%s err=%v source=%s", video.VideoId, err, videoSource)
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
	info, err := ytdl.FetchInfoWithPlatform(ctx, fullURL, videoSource)
	if err != nil {
		log.Printf("yt.info fetch_failed input=%s err=%v", idOrUrl, err)
		return nil, errcode.ErrDLYTUpstream
	}

	// 统一下载缩略图到本地
	thumbURL := info.ThumbnailUrl
	if strings.TrimSpace(thumbURL) != "" {
		if localURL, upErr := saveURLToLocal(ctx, thumbURL, buildThumbKey(info.Id, thumbURL)); upErr == nil {
			thumbURL = localURL
			log.Printf("yt.info thumb_saved video_id=%s url=%s source=%s", info.Id, thumbURL, videoSource)
		} else {
			log.Printf("yt.info thumb_save_failed video_id=%s err=%v source=%s", info.Id, upErr, videoSource)
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

	// 如果是本地静态文件，统一转换为 base64
	finalThumbURL := thumbURL
	if isLocalStaticURL(thumbURL) {
		if base64URL, err := convertLocalImageToBase64(thumbURL); err == nil {
			finalThumbURL = base64URL
			log.Printf("yt.info thumb_converted_to_base64 video_id=%s source=%s", info.Id, videoSource)
		} else {
			log.Printf("yt.info thumb_convert_base64_failed video_id=%s err=%v source=%s", info.Id, err, videoSource)
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
	return s.AudioWithPlatform(ctx, idOrUrl, "")
}

func (s *LocalYtSvc) AudioWithPlatform(ctx context.Context, idOrUrl, platform string) (*AudioResp, error) {
	// 优先使用前端传递的平台类型，否则自动检测
	videoSource := platform
	if videoSource == "" {
		videoSource = detectVideoSource(idOrUrl)
	}
	log.Printf("yt.audio video_source=%s (from_frontend=%t) input=%s", videoSource, platform != "", idOrUrl)

	// 使用新的extractVideoID函数
	videoId := extractVideoID(idOrUrl, videoSource)
	if videoId == "" {
		videoId = idOrUrl
	}

	// 构造完整URL供yt-dlp使用
	fullURL := constructFullURL(idOrUrl, videoSource)

	// 确定数据库查询的source_site
	sourceSite := videoSource
	if sourceSite == "unknown" {
		sourceSite = "youtube" // 向后兼容
	}

	// 从 dlyt 包级配置读取（由 config.InitSvc 注入），用于控制 B站音频模式
	biliMode := strings.ToLower(strings.TrimSpace(pkgOptions.BilibiliAudioMode))

	// DB 命中音频直链则直接返回；如非本地静态链接则使用 yt-dlp 直接下载到本地（在 B站 URL 模式下跳过镜像）
	var video model.YoutubeVideo
	if err := db.WithContext(ctx).Where("source_site = ? AND video_id = ?", sourceSite, videoId).First(&video).Error(); err == nil {
		if strings.TrimSpace(video.AudioUrl) != "" {
			log.Printf("yt.audio db_hit video_id=%s current_url=%s is_local=%v", video.VideoId, video.AudioUrl, isLocalStaticURL(video.AudioUrl))
			// 若为 B站且处于 URL 模式
			if videoSource == "bilibili" && biliMode == "url" {
				// 如果 DB 中存的是本地静态文件，则忽略并实时获取直链返回
				if isLocalStaticURL(video.AudioUrl) {
					bestURL, gerr := ytdl.GetBestAudioURLWithPlatform(ctx, fullURL, videoSource)
					if gerr == nil && strings.TrimSpace(bestURL) != "" {
						log.Printf("yt.audio bili_url_mode_db_local_override video_id=%s url=%s", video.VideoId, bestURL)
						return &AudioResp{Id: video.VideoId, Title: video.Title, AudioUrl: bestURL}, nil
					}
					log.Printf("yt.audio bili_url_mode_db_local_override_failed fallback_db video_id=%s err=%v", video.VideoId, gerr)
				}
				// DB 中不是本地静态（可能为远程直链），则直接返回
				return &AudioResp{Id: video.VideoId, Title: video.Title, AudioUrl: video.AudioUrl}, nil
			}
			// 否则保持原有镜像逻辑
			if !isLocalStaticURL(video.AudioUrl) {
				outBase := filepath.Join("public", "yt", "audio", video.VideoId)
				if localFile, upErr := ytdl.DownloadAudioToWithPlatform(ctx, fullURL, outBase, videoSource); upErr == nil {
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

	// B站 URL 模式：直接返回实时音频直链
	if videoSource == "bilibili" && biliMode == "url" {
		bestURL, gerr := ytdl.GetBestAudioURLWithPlatform(ctx, fullURL, videoSource)
		if gerr == nil && strings.TrimSpace(bestURL) != "" {
			log.Printf("yt.audio bili_url_mode video_id=%s url=%s", videoId, bestURL)
			return &AudioResp{Id: videoId, Title: video.Title, AudioUrl: bestURL}, nil
		}
		log.Printf("yt.audio bili_url_mode_failed fallback_local video_id=%s err=%v", videoId, gerr)
	}

	// 未命中则直接下载到本地并返回静态路径
	outBase := filepath.Join("public", "yt", "audio", videoId)
	localFile, err := ytdl.DownloadAudioToWithPlatform(ctx, fullURL, outBase, videoSource)
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
		_ = db.WithContext(ctx).Create(&model.YoutubeVideo{SourceSite: sourceSite, VideoId: videoId, AudioUrl: finalAudioURL})
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
