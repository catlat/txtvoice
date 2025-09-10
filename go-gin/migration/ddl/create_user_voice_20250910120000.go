package ddl

import (
	"go-gin/internal/migration"
)

func init() {
	migration.RegisterDDL(&CreateUserVoice20250910120000{})
}

// CreateUserVoice20250910120000 创建用户音色ID表迁移
type CreateUserVoice20250910120000 struct{}

// Up 执行迁移
func (m *CreateUserVoice20250910120000) Up(migrator *migration.DDLMigrator) error {
	return migrator.Exec(`
		CREATE TABLE user_voice (
			id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
			mobile VARCHAR(20) NOT NULL DEFAULT '' COMMENT '手机号',
			voice_id VARCHAR(50) NOT NULL DEFAULT '' COMMENT '音色ID',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
			PRIMARY KEY (id),
			UNIQUE KEY uk_mobile (mobile)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户音色ID表'
	`)
}
