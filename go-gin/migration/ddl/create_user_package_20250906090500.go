package ddl

import (
	"go-gin/internal/migration"
)

func init() {
	migration.RegisterDDL(&CreateUserPackage20250906090500{})
}

// CreateUserPackage20250906090500 创建 user_package 表
type CreateUserPackage20250906090500 struct{}

// Up 执行迁移
func (m *CreateUserPackage20250906090500) Up(migrator *migration.DDLMigrator) error {
	return migrator.Exec(`
		CREATE TABLE IF NOT EXISTS user_package (
			id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			user_identity VARCHAR(128) NOT NULL COMMENT '邮箱或手机号',
			package_id BIGINT NOT NULL COMMENT 'package.id',
			remain_asr_chars INT NOT NULL DEFAULT 0 COMMENT '剩余ASR字符',
			remain_tts_chars INT NOT NULL DEFAULT 0 COMMENT '剩余TTS字符',
			expire_at DATETIME NULL COMMENT '过期时间',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
			KEY idx_user_identity (user_identity),
			KEY idx_package_id (package_id),
			KEY idx_user_pkg (user_identity, package_id),
			CONSTRAINT fk_user_package_pkg FOREIGN KEY (package_id) REFERENCES package(id) ON DELETE RESTRICT ON UPDATE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`)
}

