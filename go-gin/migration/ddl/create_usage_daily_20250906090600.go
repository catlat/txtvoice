package ddl

import (
	"go-gin/internal/migration"
)

func init() {
	migration.RegisterDDL(&CreateUsageDaily20250906090600{})
}

// CreateUsageDaily20250906090600 创建 usage_daily 表
type CreateUsageDaily20250906090600 struct{}

// Up 执行迁移
func (m *CreateUsageDaily20250906090600) Up(migrator *migration.DDLMigrator) error {
	return migrator.Exec(`
		CREATE TABLE IF NOT EXISTS usage_daily (
			id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			user_identity VARCHAR(128) NOT NULL COMMENT '邮箱或手机号',
			date DATE NOT NULL COMMENT '日期',
			asr_chars INT NOT NULL DEFAULT 0 COMMENT '当日ASR字符',
			tts_chars INT NOT NULL DEFAULT 0 COMMENT '当日TTS字符',
			requests INT NOT NULL DEFAULT 0 COMMENT '当日请求数',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
			UNIQUE KEY uk_user_date (user_identity, date)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`)
}

