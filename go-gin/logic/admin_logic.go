package logic

import (
	"context"
	"go-gin/internal/component/db"
	"time"
)

type AdminLogic struct{}

func NewAdminLogic() *AdminLogic { return &AdminLogic{} }

// SeedQuota 将一批账号注册为正式账号并分配额度，有效期一个月（30天）
func (l *AdminLogic) SeedQuota(ctx context.Context, identities []string, asrChars, ttsChars int) error {
	if len(identities) == 0 {
		return nil
	}
	// 1) upsert account_user（正式账号）
	for _, id := range identities {
		itype := 1
		for _, ch := range id {
			if ch >= '0' && ch <= '9' {
				itype = 2
				break
			}
		}
		if err := db.WithContext(ctx).Exec(`
			INSERT INTO account_user (identity, identity_type, display_name, status, created_at, updated_at)
			VALUES(?, ?, '', 1, NOW(), NOW())
			ON DUPLICATE KEY UPDATE updated_at=NOW()
		`, id, itype).Error(); err != nil {
			return err
		}
	}
	// 2) ensure package (fixed name)
	if err := db.WithContext(ctx).Exec(`
        INSERT INTO package (name, quota_asr_chars, quota_tts_chars, monthly_reset, created_at, updated_at)
        VALUES('beta_seed', ?, ?, 0, NOW(), NOW())
        ON DUPLICATE KEY UPDATE quota_asr_chars=VALUES(quota_asr_chars), quota_tts_chars=VALUES(quota_tts_chars), updated_at=NOW()
    `, asrChars, ttsChars).Error(); err != nil {
		return err
	}

	// 取 pkg_id
	var pkgId int64
	if err := db.WithContext(ctx).Raw(`SELECT id FROM package WHERE name='beta_seed' LIMIT 1`).Scan(&pkgId).Error(); err != nil {
		return err
	}

	// 3) grant to users (replace)
	expireAt := time.Now().Add(30 * 24 * time.Hour).Format("2006-01-02 15:04:05")
	for _, id := range identities {
		_ = db.WithContext(ctx).Exec(`DELETE FROM user_package WHERE user_identity=?`, id).Error()
		if err := db.WithContext(ctx).Exec(`
            INSERT INTO user_package (user_identity, package_id, remain_asr_chars, remain_tts_chars, expire_at, created_at, updated_at)
            VALUES(?, ?, ?, ?, ?, NOW(), NOW())
        `, id, pkgId, asrChars, ttsChars, expireAt).Error(); err != nil {
			return err
		}
	}
	return nil
}
