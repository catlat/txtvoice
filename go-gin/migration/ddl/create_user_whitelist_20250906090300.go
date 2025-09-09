package ddl

import (
	"go-gin/internal/migration"
)

func init() {
	migration.RegisterDDL(&CreateUserWhitelist20250906090300{})
}

// CreateUserWhitelist20250906090300 创建 user_whitelist 表（手机号/邮箱）
type CreateUserWhitelist20250906090300 struct{}

// Up 执行迁移
func (m *CreateUserWhitelist20250906090300) Up(migrator *migration.DDLMigrator) error {
	return migrator.Exec(`
		CREATE TABLE IF NOT EXISTS user_whitelist (
			id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			identity VARCHAR(128) NOT NULL COMMENT '邮箱或手机号',
			identity_type TINYINT NOT NULL DEFAULT 0 COMMENT '0未知 1邮箱 2手机号',
			is_active TINYINT NOT NULL DEFAULT 1 COMMENT '是否启用',
			note VARCHAR(255) DEFAULT '' COMMENT '备注',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
			UNIQUE KEY uk_identity (identity)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`)
}


