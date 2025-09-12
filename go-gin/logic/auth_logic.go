package logic

import (
	"context"
	"errors"
	"log"
	"strings"

	"go-gin/const/errcode"
	"go-gin/internal/component/db"
	"go-gin/internal/security"
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

// LoginWithPassword 仅允许已存在的手机号账号使用密码登录
func (l *AuthLogic) LoginWithPassword(ctx context.Context, phone string, password string) (typing.LoginSimpleReply, error) {
	masked := func(s string) string {
		if s == "" {
			return ""
		}
		if len(s) <= 4 {
			return s
		}
		return s[:3] + "****" + s[len(s)-2:]
	}
	normalizeDigits := func(s string) string {
		var b strings.Builder
		for _, ch := range s {
			if ch >= '0' && ch <= '9' {
				b.WriteRune(ch)
			}
		}
		return b.String()
	}
	isAllDigits := func(s string) bool {
		for _, ch := range s {
			if ch < '0' || ch > '9' {
				return false
			}
		}
		return true
	}

	// 打印当前数据库名辅助排查是否连错库
	var dbName string
	if e := db.WithContext(ctx).Raw("SELECT DATABASE() AS db").Scan(&dbName).Error; e == nil {
		log.Printf("[Auth] Current DB: %s", dbName)
	}

	raw := phone
	phone = strings.TrimSpace(phone)
	password = strings.TrimSpace(password)
	norm := normalizeDigits(phone)
	log.Printf("[Auth] Login attempt phone=%s rawLen=%d trimLen=%d norm=%s normLen=%d passLen=%d allDigits=%v",
		masked(phone), len(raw), len(phone), masked(norm), len(norm), len(password), isAllDigits(phone))
	if phone == "" || password == "" {
		log.Printf("[Auth] Login reject: empty phone or password")
		return typing.LoginSimpleReply{}, errcode.ErrUserNameOrPwdFaild
	}
	// 先按 identity 精确查询，记录真实存在性
	var au model.AccountUser
	if err := db.WithContext(ctx).Debug().Where("identity=?", phone).First(&au).Error; err != nil {
		log.Printf("[Auth] No record by identity: phone=%s, try normalized match", masked(phone))
		// 尝试按去除+/-/空格后的等值或后缀匹配（处理 86 前缀场景）
		if err2 := db.WithContext(ctx).Debug().Where(
			"REPLACE(REPLACE(REPLACE(identity, '+',''),'-',''),' ','') = ? OR REPLACE(REPLACE(REPLACE(identity, '+',''),'-',''),' ','') LIKE ?",
			norm, "%"+norm,
		).First(&au).Error; err2 != nil {
			log.Printf("[Auth] Normalize match failed for phone=%s", masked(phone))
			return typing.LoginSimpleReply{}, errcode.ErrUserNameOrPwdFaild
		}
		log.Printf("[Auth] Normalize match hit: id=%d identity=%s", au.Id, masked(au.Identity))
	}
	log.Printf("[Auth] Record by identity: id=%d identityType=%d status=%d hashEmpty=%v hashLen=%d", au.Id, au.IdentityType, au.Status, au.PasswordHash == "", len(au.PasswordHash))
	// 再校验类型/状态
	if au.IdentityType != 2 || au.Status != 1 {
		log.Printf("[Auth] Reject by type/status: id=%d identityType=%d status=%d", au.Id, au.IdentityType, au.Status)
		return typing.LoginSimpleReply{}, errcode.ErrUserNameOrPwdFaild
	}
	// 若旧数据未回填哈希：允许密码==手机号的首次登录，并即时回填哈希
	if au.PasswordHash == "" {
		if password == norm || password == phone {
			log.Printf("[Auth] First-login fallback: phone==password, will backfill hash for id=%d", au.Id)
			if hash, err := security.HashPassword(password); err == nil {
				if err := db.WithContext(ctx).Model(&model.AccountUser{}).Where("id=?", au.Id).Update("password_hash", hash).Error; err != nil {
					log.Printf("[Auth] Backfill hash failed for id=%d err=%v", au.Id, err)
					return typing.LoginSimpleReply{}, errcode.ErrUserNameOrPwdFaild
				}
				au.PasswordHash = hash
				log.Printf("[Auth] Backfill hash success for id=%d", au.Id)
			} else {
				log.Printf("[Auth] HashPassword failed: id=%d err=%v", au.Id, err)
				return typing.LoginSimpleReply{}, errcode.ErrUserNameOrPwdFaild
			}
		} else {
			log.Printf("[Auth] Empty hash but password!=phone, reject id=%d", au.Id)
			return typing.LoginSimpleReply{}, errcode.ErrUserNameOrPwdFaild
		}
	}
	if !security.CheckPassword(au.PasswordHash, password) {
		log.Printf("[Auth] Password compare failed for id=%d", au.Id)
		return typing.LoginSimpleReply{}, errcode.ErrUserNameOrPwdFaild
	}
	t := token.TokenId()
	if err := token.Set(ctx, t, "identity", au.Identity); err != nil {
		log.Printf("[Auth] Token set failed for id=%d err=%v", au.Id, err)
		return typing.LoginSimpleReply{}, err
	}
	log.Printf("[Auth] Login success id=%d identity=%s tokenLen=%d", au.Id, masked(au.Identity), len(t))
	return typing.LoginSimpleReply{
		Token: t,
		User:  typing.LoginUser{Identity: au.Identity, DisplayName: au.DisplayName, Status: au.Status},
	}, nil
}

// ChangePassword 修改当前登录用户密码
func (l *AuthLogic) ChangePassword(ctx context.Context, identity string, newPassword string) error {
	masked := func(s string) string {
		if s == "" {
			return ""
		}
		if len(s) <= 4 {
			return s
		}
		return s[:3] + "****" + s[len(s)-2:]
	}
	log.Printf("[Auth] ChangePassword attempt identity=%s newLen=%d", masked(identity), len(newPassword))
	if identity == "" || len(newPassword) < 6 {
		log.Printf("[Auth] ChangePassword reject: invalid params")
		return errors.New("invalid params")
	}
	hash, err := security.HashPassword(newPassword)
	if err != nil {
		log.Printf("[Auth] HashPassword failed: identity=%s err=%v", masked(identity), err)
		return err
	}
	if err := db.WithContext(ctx).
		Model(&model.AccountUser{}).
		Where("identity=? AND status=1", identity).
		Update("password_hash", hash).Error; err != nil {
		log.Printf("[Auth] Update password failed identity=%s err=%v", masked(identity), err)
		return err
	}
	log.Printf("[Auth] ChangePassword success identity=%s", masked(identity))
	return nil
}
