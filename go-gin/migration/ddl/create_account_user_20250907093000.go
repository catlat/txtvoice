package ddl

import (
	"go-gin/internal/migration"
)

func init() {
	migration.RegisterDDL(&CreateAccountUser20250907093000{})
}

// CreateAccountUser20250907093000 创建正式账号表 account_user
type CreateAccountUser20250907093000 struct{}

// Up 执行迁移
func (m *CreateAccountUser20250907093000) Up(migrator *migration.DDLMigrator) error {
	return migrator.Exec(`
		CREATE TABLE IF NOT EXISTS account_user (
			id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			identity VARCHAR(128) NOT NULL COMMENT '邮箱或手机号',
			identity_type TINYINT NOT NULL DEFAULT 0 COMMENT '1邮箱 2手机号',
			display_name VARCHAR(128) DEFAULT '' COMMENT '展示名',
			status TINYINT NOT NULL DEFAULT 1 COMMENT '1正常 2禁用',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
			UNIQUE KEY uk_identity (identity)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`)
}
