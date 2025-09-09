package ddl

import (
	"go-gin/internal/migration"
)

func init() {
	migration.RegisterDDL(&CreateTTSHistory20250906090200{})
}

// CreateTTSHistory20250906090200 创建 tts_history 表
type CreateTTSHistory20250906090200 struct{}

// Up 执行迁移
func (m *CreateTTSHistory20250906090200) Up(migrator *migration.DDLMigrator) error {
	return migrator.Exec(`
		CREATE TABLE IF NOT EXISTS tts_history (
			id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			user_identity VARCHAR(64) NOT NULL DEFAULT '' COMMENT '用户标识',
			text_hash CHAR(64) NOT NULL COMMENT 'text+speaker+defaults 的sha256',
			text_preview VARCHAR(255) DEFAULT '' COMMENT '文本预览',
			char_count INT DEFAULT 0 COMMENT '字符数',
			speaker VARCHAR(64) NOT NULL COMMENT '说话人',
			audio_url VARCHAR(512) DEFAULT '' COMMENT '音频直链',
			request_id VARCHAR(64) DEFAULT '' COMMENT '请求ID/trace',
			status TINYINT DEFAULT 0 COMMENT '状态',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
			UNIQUE KEY uk_identity_text_speaker (user_identity, text_hash, speaker),
			KEY idx_identity_created (user_identity, created_at),
			KEY idx_created_at (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`)
}
