package logic

import (
	"context"
	"go-gin/internal/component/db"
	"go-gin/model"
)

type HistoryLogic struct{}

func NewHistoryLogic() *HistoryLogic { return &HistoryLogic{} }

func (l *HistoryLogic) ListVideos(ctx context.Context, page, size int) ([]model.YoutubeVideo, int64, error) {
	var items []model.YoutubeVideo
	var total int64

	// 构建基础查询条件
	baseQuery := db.WithContext(ctx).Model(&model.YoutubeVideo{})

	// 获取总数
	_ = baseQuery.Count(&total)

	// 获取分页数据
	q := baseQuery.Order("id desc").Limit(size).Offset((page - 1) * size)
	_ = q.Find(&items)

	return items, total, nil
}

func (l *HistoryLogic) GetVideoDetail(ctx context.Context, sourceSite, videoId string) (map[string]any, error) {
	var video model.YoutubeVideo
	_ = db.WithContext(ctx).Where("source_site=? and video_id=?", sourceSite, videoId).First(&video)

	var transcript model.YoutubeTranscript
	if video.Id != 0 {
		_ = db.WithContext(ctx).Where("video_id=?", video.Id).Order("id desc").First(&transcript)
	}

	detail := map[string]any{
		"id":            video.Id,
		"source_site":   video.SourceSite,
		"video_id":      video.VideoId,
		"title":         video.Title,
		"channel_title": video.ChannelTitle,
		"duration_sec":  video.DurationSec,
		"published_at":  video.PublishedAt,
		"thumbnail_url": video.ThumbnailUrl,
		"audio_url":     video.AudioUrl,
		"status":        video.Status,
	}
	if transcript.Id != 0 {
		detail["original_text"] = transcript.OriginalText
		detail["translated_text"] = transcript.TranslatedText
		detail["asr_char_count"] = transcript.AsrCharCount
		detail["translate_char_count"] = transcript.TranslateCharCount
		detail["utterances"] = []any{}
	}
	return detail, nil
}

func (l *HistoryLogic) ListTTS(ctx context.Context, page, size int, identity string) ([]model.TTSHistory, int64, error) {
	var items []model.TTSHistory
	var total int64

	// 构建基础查询条件
	baseQuery := db.WithContext(ctx).Model(&model.TTSHistory{})
	if identity != "" {
		baseQuery = baseQuery.Where("user_identity=?", identity)
	}

	// 获取总数
	_ = baseQuery.Count(&total)

	// 获取分页数据
	q := baseQuery.Order("id desc").Limit(size).Offset((page - 1) * size)
	_ = q.Find(&items)

	return items, total, nil
}
