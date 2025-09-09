package ddl

import (
	"go-gin/internal/migration"
)

func init() {
	migration.RegisterDDL(&CreatePackage20250906090400{})
}

// CreatePackage20250906090400 创建 package 表
type CreatePackage20250906090400 struct{}

// Up 执行迁移
func (m *CreatePackage20250906090400) Up(migrator *migration.DDLMigrator) error {
	return migrator.Exec(`
		CREATE TABLE IF NOT EXISTS package (
			id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(64) NOT NULL COMMENT '套餐名称',
			quota_asr_chars INT NOT NULL DEFAULT 0 COMMENT 'ASR可用字符数',
			quota_tts_chars INT NOT NULL DEFAULT 0 COMMENT 'TTS可用字符数',
			monthly_reset TINYINT NOT NULL DEFAULT 0 COMMENT '是否按月重置',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
			UNIQUE KEY uk_name (name)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`)
}

