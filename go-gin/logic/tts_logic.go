package logic

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"go-gin/internal/component/db"
	"go-gin/internal/component/logx"
	"go-gin/internal/metrics"
	"go-gin/model"
	"go-gin/rest/tts"
	"strings"
)

type TTSLogic struct{}

func NewTTSLogic() *TTSLogic { return &TTSLogic{} }

func (l *TTSLogic) Synthesize(ctx context.Context, identity, text, speaker string) (*model.TTSHistory, error) {
	// 强制控制台日志调试
	fmt.Printf("TTS Synthesize START: identity=%s text=%s speaker=%s\n", identity, text, speaker)

	// 幂等：sha256(identity+text+speaker)
	h := sha256.Sum256([]byte(identity + "|" + text + "|" + speaker))
	textHash := hex.EncodeToString(h[:])
	fmt.Printf("TTS textHash=%s\n", textHash)

	var item model.TTSHistory
	dbResult := db.WithContext(ctx).Where("user_identity=? AND text_hash=? AND speaker=?", identity, textHash, speaker).First(&item)
	fmt.Printf("TTS DB query result: error=%v, found_id=%d\n", dbResult.Error(), item.Id)

	if item.Id != 0 {
		fmt.Printf("TTS history HIT existing record id=%d\n", item.Id)
		logx.WithContext(ctx).Info("tts_history_hit", map[string]any{"id": item.Id, "identity": identity, "speaker": speaker})
		return &item, nil
	}

	// 外部 TTS
	fmt.Printf("TTS calling external service...\n")
	resp, err := tts.Svc.Synthesize(ctx, text, speaker)
	if err != nil {
		fmt.Printf("TTS external service failed: %v\n", err)
		return nil, err
	}
	fmt.Printf("TTS external service OK\n")

	// 生成可用的 audio_url（若服务未返回直链，则使用内联 dataURL）
	audioURL := resp.AudioUrl
	if audioURL == "" && len(resp.Audio) > 0 {
		encoded := base64.StdEncoding.EncodeToString(resp.Audio)
		audioURL = "data:audio/mp3;base64," + encoded
	}
	fmt.Printf("TTS audioURL type: %s, len=%d\n", func() string {
		if strings.HasPrefix(audioURL, "data:") {
			return "data"
		}
		return "remote"
	}(), len(audioURL))

	// 入库
	preview := text
	r := []rune(preview)
	if len(r) > 255 {
		preview = string(r[:255])
	}
	item = model.TTSHistory{
		UserIdentity: identity,
		TextHash:     textHash,
		TextPreview:  preview,
		CharCount:    len([]rune(text)),
		Speaker:      speaker,
		AudioUrl:     audioURL,
		RequestId:    resp.RequestId,
		Status:       0,
	}
	fmt.Printf("TTS creating DB record: identity=%s, textHash=%s, preview=%s\n", item.UserIdentity, item.TextHash, item.TextPreview)

	if err := db.WithContext(ctx).Create(&item).Error(); err != nil {
		fmt.Printf("TTS DB create FAILED: %v\n", err)
		logx.WithContext(ctx).Error("tts_history_create_failed", map[string]any{"identity": identity, "speaker": speaker, "err": err.Error()})
		// 不中断主流程，仍返回音频结果，但提示日志排查 DB/Migration
		return &item, nil
	}

	fmt.Printf("TTS DB create SUCCESS: new_id=%d\n", item.Id)
	logx.WithContext(ctx).Info("tts_history_created", map[string]any{"id": item.Id, "identity": identity, "speaker": speaker})

	_ = metrics.AddUsage(ctx, identity, 0, item.CharCount, 1)
	return &item, nil
}
