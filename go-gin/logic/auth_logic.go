package logic

import (
	"context"
	"go-gin/internal/component/db"
	"go-gin/internal/token"
	"go-gin/model"
	"go-gin/typing"
)

type AuthLogic struct{}

func NewAuthLogic() *AuthLogic { return &AuthLogic{} }

// LoginSimple 懒注册到 account_user，并发放会话 token
func (l *AuthLogic) LoginSimple(ctx context.Context, identity string) (typing.LoginSimpleReply, error) {
	// 简单判断身份类型：全数字视为手机，否则邮箱
	itype := 1
	for _, ch := range identity {
		if ch >= '0' && ch <= '9' {
			itype = 2
			break
		}
	}
	// upsert account_user
	if err := db.WithContext(ctx).Exec(`
		INSERT INTO account_user (identity, identity_type, display_name, status, created_at, updated_at)
		VALUES(?, ?, '', 1, NOW(), NOW())
		ON DUPLICATE KEY UPDATE updated_at=NOW()
	`, identity, itype).Error(); err != nil {
		return typing.LoginSimpleReply{}, err
	}
	// 查询用户信息
	var au model.AccountUser
	_ = db.WithContext(ctx).Where("identity=?", identity).First(&au)
	// issue token
	t := token.TokenId()
	if err := token.Set(ctx, t, "identity", identity); err != nil {
		return typing.LoginSimpleReply{}, err
	}
	return typing.LoginSimpleReply{
		Token: t,
		User:  typing.LoginUser{Identity: au.Identity, DisplayName: au.DisplayName, Status: au.Status},
	}, nil
}

func (l *AuthLogic) Logout(ctx context.Context, tokenId string) error {
	return token.Flush(ctx, tokenId)
}
