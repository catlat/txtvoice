package metrics

import (
	"context"
	"fmt"
	"go-gin/internal/component/db"
	"time"
)

// AddUsage 将用量写入 usage_daily，按 (user_identity, date) 聚合
func AddUsage(ctx context.Context, identity string, asrChars, ttsChars, requests int) error {
	if identity == "" {
		identity = "guest"
	}

	today := time.Now().Format("2006-01-02")

	// 添加调试日志
	fmt.Printf("AddUsage called: identity=%s, date=%s, asr=%d, tts=%d, requests=%d\n",
		identity, today, asrChars, ttsChars, requests)

	sql := `INSERT INTO usage_daily (user_identity, date, asr_chars, tts_chars, requests, created_at)
            VALUES (?, ?, ?, ?, ?, NOW())
            ON DUPLICATE KEY UPDATE
                asr_chars = asr_chars + VALUES(asr_chars),
                tts_chars = tts_chars + VALUES(tts_chars),
                requests = requests + VALUES(requests)`

	err := db.WithContext(ctx).Exec(sql, identity, today, asrChars, ttsChars, requests).Error()
	if err != nil {
		fmt.Printf("AddUsage error: %v\n", err)
		return err
	}

	fmt.Printf("AddUsage success for identity=%s\n", identity)
	return nil
}
