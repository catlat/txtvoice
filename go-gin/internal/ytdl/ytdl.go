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

	"go-gin/internal/component/redisx"

	"github.com/redis/go-redis/v9"
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
	Timestamp  int64   `json:"timestamp"` // 时间戳，适用于Bilibili
	// YouTube的缩略图数组
	Thumbnails []struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"thumbnails"`
	// Bilibili的单个缩略图
	Thumbnail string `json:"thumbnail"`
}

func getBin() string {
	if b := strings.TrimSpace(os.Getenv("YTDL_BIN")); b != "" {
		return b
	}
	return "yt-dlp"
}

// getCookies 从 Redis 读取 Netscape Cookie 文本，写入临时文件并返回路径（不兼容旧版路径值）
func getCookies(ctx context.Context, platform string) string {
	if platform == "" {
		platform = "youtube" // 默认平台
	}

	redisKey := fmt.Sprintf("ytdl:cookies:%s", platform)
	redisVal, err := redisx.Client().Get(ctx, redisKey).Result()
	if err == redis.Nil {
		// key不存在，返回空字符串
		return ""
	} else if err != nil {
		log.Printf("ytdl: failed to get cookies from redis for platform %s: %v", platform, err)
		return ""
	}

	v := strings.TrimSpace(redisVal)
	if v == "" {
		return ""
	}

	// 始终按文本处理（Netscape格式）
	tmpFile, createErr := os.CreateTemp("", fmt.Sprintf("ytcookies_%s_*.txt", platform))
	if createErr != nil {
		log.Printf("ytdl: create temp cookie file failed platform=%s err=%v", platform, createErr)
		return ""
	}
	if _, writeErr := tmpFile.WriteString(v); writeErr != nil {
		_ = tmpFile.Close()
		log.Printf("ytdl: write temp cookie file failed platform=%s err=%v", platform, writeErr)
		return ""
	}
	_ = tmpFile.Close()
	log.Printf("ytdl: using cookies from redis (written to temp) platform=%s file=%s len=%d", platform, tmpFile.Name(), len(v))
	return tmpFile.Name()
}

func getProxy() string { return strings.TrimSpace(os.Getenv("YTDL_PROXY")) }

func getAudioFormat() string {
	if f := strings.TrimSpace(os.Getenv("YTDL_AUDIO_FORMAT")); f != "" {
		return f
	}
	// 优化为最低质量音频，适合语音识别
	return "worstaudio/bestaudio[abr<=32]/bestaudio[abr<=64]/bestaudio"
}

// FetchInfo 使用 yt-dlp 获取视频信息（本地执行，不依赖外部服务）
func FetchInfo(ctx context.Context, idOrURL string) (*Info, error) {
	return FetchInfoWithPlatform(ctx, idOrURL, "")
}

// FetchInfoWithPlatform 使用指定平台的cookie获取视频信息
func FetchInfoWithPlatform(ctx context.Context, idOrURL, platform string) (*Info, error) {
	bin := getBin()
	args := []string{"-J", "--no-playlist"}

	// 从Redis读取平台的cookie文件路径
	if c := getCookies(ctx, platform); c != "" {
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

	// 处理缩略图：优先使用Bilibili的thumbnail字段，否则使用thumbnails数组
	var thumbnailURL string
	if data.Thumbnail != "" {
		// Bilibili 单个缩略图
		thumbnailURL = data.Thumbnail
		log.Printf("ytdl: using single thumbnail url=%s", thumbnailURL)
	} else {
		// YouTube 缩略图数组，选择最大分辨率的
		bestArea := 0
		log.Printf("ytdl: thumbnails count=%d", len(data.Thumbnails))
		for _, t := range data.Thumbnails {
			area := t.Width * t.Height
			if area > bestArea {
				bestArea = area
				thumbnailURL = t.URL
			}
		}
		log.Printf("ytdl: chosen thumb area=%d url=%s", bestArea, thumbnailURL)
	}

	// 处理发布日期：优先使用upload_date，如果没有则使用timestamp转换
	publishDate := data.UploadDate
	if publishDate == "" && data.Timestamp > 0 {
		// 将timestamp转换为YYYYMMDD格式
		t := time.Unix(data.Timestamp, 0)
		publishDate = t.Format("20060102")
		log.Printf("ytdl: converted timestamp %d to date %s", data.Timestamp, publishDate)
	}

	return &Info{
		Id:           data.ID,
		Title:        data.Title,
		Author:       author,
		DurationSec:  int(data.Duration),
		Views:        data.ViewCount,
		PublishDate:  publishDate,
		ThumbnailUrl: thumbnailURL,
	}, nil
}

// GetBestAudioURL 使用 yt-dlp 获取音频直链
func GetBestAudioURL(ctx context.Context, idOrURL string) (string, error) {
	return GetBestAudioURLWithPlatform(ctx, idOrURL, "")
}

// GetBestAudioURLWithPlatform 使用指定平台的cookie获取音频直链
func GetBestAudioURLWithPlatform(ctx context.Context, idOrURL, platform string) (string, error) {
	bin := getBin()
	format := getAudioFormat()
	args := []string{"-f", format, "-g", "--no-playlist"}

	// 从Redis读取平台的cookie文件路径
	if c := getCookies(ctx, platform); c != "" {
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
	return DownloadAudioToWithPlatform(ctx, idOrURL, outBase, "")
}

// DownloadAudioToWithPlatform 使用指定平台的cookie下载音频到本地文件
func DownloadAudioToWithPlatform(ctx context.Context, idOrURL, outBase, platform string) (string, error) {
	bin := getBin()
	format := getAudioFormat()
	outTemplate := outBase + ".%(ext)s"
	args := []string{"-f", format, "--no-playlist", "-o", outTemplate}

	// 从Redis读取平台的cookie文件路径
	if c := getCookies(ctx, platform); c != "" {
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
