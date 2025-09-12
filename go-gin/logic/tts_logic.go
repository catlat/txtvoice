package logic

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"go-gin/const/errcode"
	"go-gin/internal/component/db"
	"go-gin/internal/component/logx"
	"go-gin/internal/metrics"
	"go-gin/model"
	"go-gin/rest/dlyt"
	"go-gin/rest/tts"
	"strings"
	"time"
)

type TTSLogic struct{}

func NewTTSLogic() *TTSLogic { return &TTSLogic{} }

func (l *TTSLogic) Synthesize(ctx context.Context, identity, text, speaker string, useMyVoice bool) (*model.TTSHistory, error) {
	// 简洁日志记录
	fmt.Printf("TTS START: identity=%s useMyVoice=%v textLen=%d\n", identity, useMyVoice, len(text))

	// 预检余额（严格：不足直接拒绝）
	need := len([]rune(text))
	if ok, err := l.hasEnoughTTSBalance(ctx, identity, need); err == nil && !ok {
		return nil, errcode.ErrQuotaNotEnough
	}

	// 决定有效 speaker 与资源（不回退）
	effectiveSpeaker := speaker
	resourceId := "volc.service_type.10029" // 大模型语音合成（字符版）
	if useMyVoice {
		uv, err := model.NewUserVoiceModel().GetByMobile(ctx, identity)
		if err != nil || uv == nil || uv.VoiceId == "" {
			logx.WithContext(ctx).Warn("user_voice_not_configured", map[string]any{"identity": identity})
			return nil, errcode.ErrUserVoiceNotConfigured
		}
		effectiveSpeaker = uv.VoiceId
		resourceId = "volc.megatts.default" // 声音复刻2.0（字符版）
		fmt.Printf("TTS using my voice: identity=%s, voice_id=%s\n", identity, uv.VoiceId)
	}

	// 幂等：sha256(identity|text|effectiveSpeaker)
	h := sha256.Sum256([]byte(identity + "|" + text + "|" + effectiveSpeaker))
	textHash := hex.EncodeToString(h[:])

	var item model.TTSHistory
	db.WithContext(ctx).Where("user_identity=? AND text_hash=? AND speaker=?", identity, textHash, effectiveSpeaker).First(&item)

	if item.Id != 0 {
		fmt.Printf("TTS cache hit: id=%d\n", item.Id)
		logx.WithContext(ctx).Info("tts_history_hit", map[string]any{"id": item.Id, "identity": identity, "speaker": effectiveSpeaker})

		// 缓存命中时也要进行余额预检（按历史记录的字符数）
		if ok, err := l.hasEnoughTTSBalance(ctx, identity, item.CharCount); err == nil && !ok {
			return nil, errcode.ErrQuotaNotEnough
		}

		// 缓存命中时也需要记录使用统计
		if err := metrics.AddUsage(ctx, identity, 0, item.CharCount, 1); err != nil {
			fmt.Printf("TTS cache hit AddUsage failed: identity=%s, chars=%d, error=%v\n", identity, item.CharCount, err)
			logx.WithContext(ctx).Error("tts_cache_usage_record_failed", map[string]any{"identity": identity, "chars": item.CharCount, "error": err.Error()})
		} else {
			fmt.Printf("TTS cache hit AddUsage success: identity=%s, chars=%d\n", identity, item.CharCount)
		}

		// 缓存命中时也需要扣减套餐余额
		if err := l.deductTTSBalance(ctx, identity, item.CharCount); err != nil {
			fmt.Printf("TTS cache hit balance deduction failed: identity=%s, chars=%d, error=%v\n", identity, item.CharCount, err)
			logx.WithContext(ctx).Error("tts_cache_balance_deduction_failed", map[string]any{"identity": identity, "chars": item.CharCount, "error": err.Error()})
		} else {
			fmt.Printf("TTS cache hit balance deducted: identity=%s, chars=%d\n", identity, item.CharCount)
		}

		return &item, nil
	}

	// 外部 TTS（按指定资源调用）
	fmt.Printf("TTS calling external service: resource=%s speaker=%s\n", resourceId, effectiveSpeaker)
	resp, err := tts.Svc.SynthesizeWithResource(ctx, text, effectiveSpeaker, resourceId)
	if err != nil {
		fmt.Printf("TTS failed: %v\n", err)
		return nil, err
	}

	// 将音频保存到七牛云，数据库仅存公网链接
	audioURL := ""
	// 构造稳定的对象键：tts/{identity}/{hash8}-{unix}.mp3
	hash8 := textHash
	if len(hash8) > 8 {
		hash8 = textHash[:8]
	}
	safeIdentity := strings.ReplaceAll(identity, "|", "_")
	key := fmt.Sprintf("tts/%s/%s-%d.mp3", safeIdentity, hash8, time.Now().Unix())

	// 优先服务端 Fetch（如果上游给了 URL），否则直接上传字节
	var upErr error
	if strings.TrimSpace(resp.AudioUrl) != "" {
		if url, err := dlyt.FetchToQiniu(ctx, key, resp.AudioUrl); err == nil {
			audioURL = url
		} else {
			upErr = err
		}
	}
	if audioURL == "" && len(resp.Audio) > 0 {
		if url, err := dlyt.UploadBytesToQiniu(ctx, key, resp.Audio, "audio/mpeg"); err == nil {
			audioURL = url
		} else {
			upErr = err
		}
	}
	// 七牛不可用或上传失败时回退为 data URL，保证不阻断主流程
	if audioURL == "" && len(resp.Audio) > 0 {
		encoded := base64.StdEncoding.EncodeToString(resp.Audio)
		audioURL = "data:audio/mp3;base64," + encoded
		if upErr != nil {
			logx.WithContext(ctx).Warn("tts_qiniu_upload_failed_fallback", map[string]any{"err": upErr.Error()})
		}
	}

	// 入库
	preview := text
	item = model.TTSHistory{
		UserIdentity: identity,
		TextHash:     textHash,
		TextPreview:  preview,
		CharCount:    len([]rune(text)),
		Speaker:      effectiveSpeaker,
		AudioUrl:     audioURL,
		RequestId:    resp.RequestId,
		Status:       0,
	}

	if err := db.WithContext(ctx).Create(&item).Error(); err != nil {
		logx.WithContext(ctx).Error("tts_history_create_failed", map[string]any{"identity": identity, "speaker": effectiveSpeaker, "err": err.Error()})
		// 不中断主流程，仍返回音频结果，但提示日志排查 DB/Migration
		return &item, nil
	}

	fmt.Printf("TTS success: saved id=%d\n", item.Id)
	logx.WithContext(ctx).Info("tts_history_created", map[string]any{"id": item.Id, "identity": identity, "speaker": effectiveSpeaker})

	// 记录使用统计 - 确保即使统计失败也不影响主流程
	if err := metrics.AddUsage(ctx, identity, 0, item.CharCount, 1); err != nil {
		fmt.Printf("TTS AddUsage failed: identity=%s, chars=%d, error=%v\n", identity, item.CharCount, err)
		logx.WithContext(ctx).Error("tts_usage_record_failed", map[string]any{"identity": identity, "chars": item.CharCount, "error": err.Error()})
	} else {
		fmt.Printf("TTS AddUsage success: identity=%s, chars=%d\n", identity, item.CharCount)
	}

	// 扣减套餐余额 - 即使扣减失败也不影响主流程
	if err := l.deductTTSBalance(ctx, identity, item.CharCount); err != nil {
		fmt.Printf("TTS balance deduction failed: identity=%s, chars=%d, error=%v\n", identity, item.CharCount, err)
		logx.WithContext(ctx).Error("tts_balance_deduction_failed", map[string]any{"identity": identity, "chars": item.CharCount, "error": err.Error()})
	} else {
		fmt.Printf("TTS balance deducted: identity=%s, chars=%d\n", identity, item.CharCount)
	}

	return &item, nil
}

// deductTTSBalance 扣减用户TTS套餐余额
func (l *TTSLogic) deductTTSBalance(ctx context.Context, identity string, chars int) error {
	if identity == "" || chars <= 0 {
		return nil
	}

	// 更新用户套餐余额，按优先级扣减（先到期的先扣）
	sql := `UPDATE user_package 
			SET remain_tts_chars = GREATEST(0, remain_tts_chars - ?),
				updated_at = NOW()
			WHERE user_identity = ? 
			AND remain_tts_chars > 0 
			AND (expire_at IS NULL OR expire_at > NOW())
			ORDER BY expire_at ASC 
			LIMIT 1`

	result := db.WithContext(ctx).Exec(sql, chars, identity)
	if result.Error() != nil {
		return result.Error()
	}

	// 如果没有更新任何记录，说明用户没有可用余额，但不报错（允许透支使用）
	if result.RowsAffected == 0 {
		fmt.Printf("TTS balance deduction: no available balance for identity=%s, chars=%d\n", identity, chars)
	}

	return nil
}

// hasEnoughTTSBalance 返回是否有足够的 TTS 余额
func (l *TTSLogic) hasEnoughTTSBalance(ctx context.Context, identity string, need int) (bool, error) {
	if identity == "" || need <= 0 {
		return true, nil
	}
	var remain int
	row := db.WithContext(ctx).Raw("SELECT COALESCE(SUM(remain_tts_chars),0) FROM user_package WHERE user_identity = ? AND (expire_at IS NULL OR expire_at > NOW())", identity).Row()
	if err := row.Err(); err != nil {
		return true, err // 查询异常时不阻断
	}
	_ = row.Scan(&remain)
	return remain >= need, nil
}
