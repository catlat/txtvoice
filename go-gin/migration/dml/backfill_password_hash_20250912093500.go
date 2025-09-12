package dml

import (
	"go-gin/internal/migration"
	"go-gin/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func init() {
	migration.RegisterDML(&BackfillPasswordHash20250912093500{})
}

// BackfillPasswordHash20250912093500 为手机号账号回填默认密码（密码=手机号）
type BackfillPasswordHash20250912093500 struct{}

func (m *BackfillPasswordHash20250912093500) Desc() string {
	return "backfill account_user.password_hash with bcrypt(identity) where identity_type=2 and hash empty"
}

func (m *BackfillPasswordHash20250912093500) Handle(db *gorm.DB) error {
	// 仅处理列存在的情况
	if !db.Migrator().HasColumn(&model.AccountUser{}, "password_hash") {
		return nil
	}
	const pageSize = 200
	offset := 0
	for {
		var users []model.AccountUser
		if err := db.Model(&model.AccountUser{}).
			Where("identity_type = ? AND (password_hash IS NULL OR password_hash = '')", 2).
			Order("id ASC").
			Limit(pageSize).
			Offset(offset).
			Find(&users).Error; err != nil {
			return err
		}
		if len(users) == 0 {
			break
		}
		for _, u := range users {
			if u.Identity == "" {
				continue
			}
			hash, err := bcrypt.GenerateFromPassword([]byte(u.Identity), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			if err := db.Model(&model.AccountUser{}).
				Where("id = ?", u.Id).
				Update("password_hash", string(hash)).Error; err != nil {
				return err
			}
		}
		if len(users) < pageSize {
			break
		}
		offset += pageSize
	}
	return nil
}
