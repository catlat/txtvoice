package dml

import (
	"go-gin/internal/migration"

	"gorm.io/gorm"
)

func init() { migration.RegisterDML(&MergeAccounts20250907100500{}) }

type MergeAccounts20250907100500 struct{}

func (m *MergeAccounts20250907100500) Desc() string {
	return "ensure user_package identities exist in account_user"
}

func (m *MergeAccounts20250907100500) Handle(db *gorm.DB) error {
	// 确保 user_package 的身份在 account_user 中存在
	if err := db.Exec(`
        INSERT INTO account_user (identity, identity_type, display_name, status, created_at, updated_at)
        SELECT up.user_identity,
               CASE WHEN up.user_identity REGEXP '^[0-9]+$' THEN 2 ELSE 1 END AS identity_type,
               '', 1, NOW(), NOW()
        FROM user_package up
        LEFT JOIN account_user au ON au.identity = up.user_identity
        WHERE au.id IS NULL;
    `).Error; err != nil {
		return err
	}
	return nil
}
