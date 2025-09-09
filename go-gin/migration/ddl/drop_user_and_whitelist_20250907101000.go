package ddl

import (
	"go-gin/internal/migration"
)

func init() { migration.RegisterDDL(&DropUserAndWhitelist20250907101000{}) }

type DropUserAndWhitelist20250907101000 struct{}

func (m *DropUserAndWhitelist20250907101000) Up(migrator *migration.DDLMigrator) error {
	// 安全删除演示/白名单表（若存在）
	_ = migrator.Exec("DROP TABLE IF EXISTS user_whitelist;")
	_ = migrator.Exec("DROP TABLE IF EXISTS user;")
	return nil
}
