package metrics

import (
	"context"
	"go-gin/internal/component/db"
	"time"
)

// AddUsage 将用量写入 usage_daily，按 (user_identity, date) 聚合
func AddUsage(ctx context.Context, identity string, asrChars, ttsChars, requests int) error {
	if identity == "" {
		identity = "guest"
	}
	today := time.Now().Format("2006-01-02")
	sql := `INSERT INTO usage_daily (user_identity, date, asr_chars, tts_chars, requests, created_at)
            VALUES (?, ?, ?, ?, ?, NOW())
            ON DUPLICATE KEY UPDATE
                asr_chars = asr_chars + VALUES(asr_chars),
                tts_chars = tts_chars + VALUES(tts_chars),
                requests = requests + VALUES(requests)`
	return db.WithContext(ctx).Exec(sql, identity, today, asrChars, ttsChars, requests).Error()
}
