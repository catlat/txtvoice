package logic

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/url"
	"strings"

	"go-gin/const/errcode"
	"go-gin/internal/component/db"
	"go-gin/internal/metrics"
	"go-gin/model"
	"go-gin/rest/asr"
	"go-gin/rest/dlyt"
	"go-gin/rest/translate"
)

type TranscriptLogic struct{}

func NewTranscriptLogic() *TranscriptLogic { return &TranscriptLogic{} }

// GetOrCreate 流程：dlyt.Info -> upsert youtube_video -> dlyt.Audio -> ASR -> Translate -> upsert youtube_transcript
func (l *TranscriptLogic) GetOrCreate(ctx context.Context, idOrUrl string, targetLang string, identity string) (*model.YoutubeTranscript, error) {
	return l.GetOrCreateWithPlatform(ctx, idOrUrl, targetLang, identity, "")
}

// GetOrCreateWithPlatform 带平台参数的转录处理方法
func (l *TranscriptLogic) GetOrCreateWithPlatform(ctx context.Context, idOrUrl string, targetLang string, identity string, platform string) (*model.YoutubeTranscript, error) {
	log.Printf("[Transcript] 开始处理转录请求 - IdOrUrl: %s, TargetLang: %s, Identity: %s, Platform: %s", idOrUrl, targetLang, identity, platform)

	if targetLang == "" {
		targetLang = "zh"
	}

	log.Printf("[Transcript] Step 1: 获取视频信息")
	info, err := dlyt.Svc.InfoWithPlatform(ctx, idOrUrl, platform)
	if err != nil {
		log.Printf("[Transcript] 获取视频信息失败 - Error: %v", err)
		return nil, err
	}
	log.Printf("[Transcript] 视频信息获取成功 - VideoId: %s, Title: %s, Duration: %d", info.Id, info.Title, info.DurationSec)

	// 根据平台确定source_site
	sourceSite := platform
	if sourceSite == "" {
		sourceSite = "youtube" // 向后兼容
	}

	log.Printf("[Transcript] Step 2: 检查/创建视频记录")
	var video model.YoutubeVideo
	if err := db.WithContext(ctx).Where("source_site = ? AND video_id = ?", sourceSite, info.Id).First(&video).Error(); err != nil {
		log.Printf("[Transcript] 创建新视频记录 - VideoId: %s, SourceSite: %s", info.Id, sourceSite)
		video = model.YoutubeVideo{
			SourceSite:   sourceSite,
			VideoId:      info.Id,
			Title:        info.Title,
			ChannelTitle: info.Author,
			DurationSec:  info.DurationSec,
			ThumbnailUrl: info.ThumbnailUrl,
		}
		_ = db.WithContext(ctx).Create(&video)
	} else {
		log.Printf("[Transcript] 使用已存在的视频记录 - VideoId: %s, DBId: %d", info.Id, video.Id)
	}

	log.Printf("[Transcript] Step 3: 获取音频URL")
	audio, err := dlyt.Svc.AudioWithPlatform(ctx, idOrUrl, platform)
	if err != nil {
		log.Printf("[Transcript] 获取音频失败 - Error: %v", err)
		return nil, err
	}
	log.Printf("[Transcript] 音频获取成功 - AudioUrl: %s", audio.AudioUrl)

	// 验证音频URL有效性
	if err := validateAudioUrl(audio.AudioUrl); err != nil {
		log.Printf("[Transcript] 音频URL无效 - Error: %v", err)
		return nil, err
	}

	// 检查音频来源类型
	if strings.HasPrefix(audio.AudioUrl, "/static/") {
		log.Printf("[Transcript] 使用本地静态音频文件 - Path: %s", audio.AudioUrl)
	} else if isQiniuUrl(audio.AudioUrl) {
		log.Printf("[Transcript] 使用七牛存储音频文件 - URL: %s", audio.AudioUrl)
	} else {
		log.Printf("[Transcript] 使用外部音频URL - URL: %s", audio.AudioUrl)
	}

	log.Printf("[Transcript] Step 4: 执行语音识别")
	asrResp, err := asr.Svc.Recognize(ctx, audio.AudioUrl)
	if err != nil {
		log.Printf("[Transcript] 语音识别失败 - Error: %v", err)
		return nil, err
	}
	log.Printf("[Transcript] 语音识别成功 - CharCount: %d, TextPreview: %.100s...",
		asrResp.CharCount, asrResp.Text)

	// 根据平台决定是否翻译
	var finalText string
	var translateCharCount int

	if sourceSite == "bilibili" {
		// Bilibili视频直接使用ASR结果，不翻译
		log.Printf("[Transcript] Step 5: Bilibili视频跳过翻译，直接使用ASR结果")
		finalText = asrResp.Text
		translateCharCount = 0
	} else {
		// YouTube视频需要翻译
		log.Printf("[Transcript] Step 5: 执行翻译")
		trResp, err := translate.Svc.TranslateToZh(ctx, asrResp.Text)
		if err != nil {
			log.Printf("[Transcript] 翻译失败 - Error: %v", err)
			return nil, err
		}
		log.Printf("[Transcript] 翻译成功 - CharCount: %d, TextPreview: %.100s...",
			trResp.CharCount, trResp.Text)
		finalText = trResp.Text
		translateCharCount = trResp.CharCount
	}

	log.Printf("[Transcript] Step 6: 保存转录结果")
	var transcript model.YoutubeTranscript
	if err := db.WithContext(ctx).Where("video_id = ? AND language = ?", video.Id, targetLang).First(&transcript).Error(); err != nil {
		log.Printf("[Transcript] 创建新转录记录 - VideoId: %d, Language: %s", video.Id, targetLang)
		transcript = model.YoutubeTranscript{
			VideoId:            video.Id,
			Language:           targetLang,
			OriginalText:       asrResp.Text,
			TranslatedText:     finalText,
			AsrCharCount:       asrResp.CharCount,
			TranslateCharCount: translateCharCount,
		}
		_ = db.WithContext(ctx).Create(&transcript)
	} else {
		log.Printf("[Transcript] 更新已存在的转录记录 - TranscriptId: %d", transcript.Id)
		updates := map[string]any{
			"original_text":        asrResp.Text,
			"translated_text":      finalText,
			"asr_char_count":       asrResp.CharCount,
			"translate_char_count": translateCharCount,
		}
		_ = db.WithContext(ctx).Model(&model.YoutubeTranscript{}).Where("id = ?", transcript.Id).Updates(updates)
	}

	log.Printf("[Transcript] Step 7: 更新使用统计")
	_ = metrics.AddUsage(ctx, identity, asrResp.CharCount, 0, 1)
	if translateCharCount > 0 {
		_ = metrics.AddUsage(ctx, identity, 0, translateCharCount, 0)
	}

	log.Printf("[Transcript] Step 8: 扣减套餐余额")
	// 按最终返回给前端的中文字符数扣费（finalText）
	finalCharCount := len([]rune(finalText))
	if err := l.deductASRBalance(ctx, identity, finalCharCount); err != nil {
		log.Printf("[Transcript] ASR balance deduction failed: identity=%s, chars=%d, error=%v", identity, finalCharCount, err)
	} else {
		log.Printf("[Transcript] ASR balance deducted: identity=%s, chars=%d (final text)", identity, finalCharCount)
	}

	log.Printf("[Transcript] 转录处理完成 - TranscriptId: %d, ASR字符数: %d, 翻译字符数: %d",
		transcript.Id, asrResp.CharCount, translateCharCount)

	return &transcript, nil
}

func validateAudioUrl(audioUrl string) error {
	if audioUrl == "" {
		return errcode.ErrASRUpstream // 使用ASR错误，因为这会导致ASR失败
	}

	// 检查是否为本地静态路径
	if strings.HasPrefix(audioUrl, "/static/") {
		log.Printf("[Transcript] 检测到本地静态音频路径 - Path: %s", audioUrl)
		// 检查是否为音频文件扩展名
		path := strings.ToLower(audioUrl)
		audioExts := []string{".mp3", ".wav", ".m4a", ".aac", ".ogg", ".flac"}
		hasAudioExt := false
		for _, ext := range audioExts {
			if strings.HasSuffix(path, ext) {
				hasAudioExt = true
				break
			}
		}
		if !hasAudioExt {
			log.Printf("[Transcript] 警告: 本地静态路径可能不是音频文件 - Path: %s", audioUrl)
		}
		return nil
	}

	// 检查URL格式（用于外部URL）
	parsedUrl, err := url.Parse(audioUrl)
	if err != nil {
		log.Printf("[Transcript] 音频URL解析失败 - URL: %s, Error: %v", audioUrl, err)
		return errcode.ErrASRUpstream
	}

	if parsedUrl.Scheme == "" || parsedUrl.Host == "" {
		log.Printf("[Transcript] 音频URL格式无效 - URL: %s", audioUrl)
		return errcode.ErrASRUpstream
	}

	// 检查是否为可能的音频文件
	path := strings.ToLower(parsedUrl.Path)
	audioExts := []string{".mp3", ".wav", ".m4a", ".aac", ".ogg", ".flac"}
	hasAudioExt := false
	for _, ext := range audioExts {
		if strings.HasSuffix(path, ext) {
			hasAudioExt = true
			break
		}
	}

	if !hasAudioExt {
		log.Printf("[Transcript] 警告: 音频URL可能不是音频文件 - URL: %s", audioUrl)
		// 不返回错误，因为有些服务可能没有文件扩展名
	}

	return nil
}

func isQiniuUrl(u string) bool {
	if u == "" {
		return false
	}
	// 检查是否包含七牛域名（oss.duckai.cn）
	return strings.Contains(u, "oss.duckai.cn") ||
		strings.Contains(u, "qiniu") ||
		strings.Contains(u, "qn.") ||
		strings.Contains(u, ".qiniudn.com") ||
		strings.Contains(u, ".clouddn.com")
}

func HashTextSpeaker(text, speaker string) string {
	h := sha256.Sum256([]byte(text + "|" + speaker))
	return hex.EncodeToString(h[:])
}

// deductASRBalance 扣减用户ASR套餐余额
func (l *TranscriptLogic) deductASRBalance(ctx context.Context, identity string, chars int) error {
	if identity == "" || chars <= 0 {
		return nil
	}

	// 更新用户套餐余额，按优先级扣减（先到期的先扣）
	sql := `UPDATE user_package 
			SET remain_asr_chars = GREATEST(0, remain_asr_chars - ?),
				updated_at = NOW()
			WHERE user_identity = ? 
			AND remain_asr_chars > 0 
			AND (expire_at IS NULL OR expire_at > NOW())
			ORDER BY expire_at ASC 
			LIMIT 1`

	result := db.WithContext(ctx).Exec(sql, chars, identity)
	if result.Error() != nil {
		return result.Error()
	}

	// 如果没有更新任何记录，说明用户没有可用余额，但不报错（允许透支使用）
	if result.RowsAffected == 0 {
		log.Printf("[Transcript] ASR balance deduction: no available balance for identity=%s, chars=%d", identity, chars)
	}

	return nil
}
