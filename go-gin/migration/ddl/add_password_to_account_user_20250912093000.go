package ddl

import (
	"go-gin/internal/migration"
)

func init() {
	migration.RegisterDDL(&AddPasswordToAccountUser20250912093000{})
}

// AddPasswordToAccountUser20250912093000 为 account_user 增加 password_hash 字段
type AddPasswordToAccountUser20250912093000 struct{}

// Up 执行迁移
func (m *AddPasswordToAccountUser20250912093000) Up(migrator *migration.DDLMigrator) error {
	// 添加 password_hash 列（若不存在）
	if !migrator.HasColumn("account_user", "password_hash") {
		if err := migrator.Exec(`
            ALTER TABLE account_user
            ADD COLUMN password_hash VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'bcrypt 哈希密码' AFTER display_name;
        `); err != nil {
			return err
		}
	}
	return nil
}
