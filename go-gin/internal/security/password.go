package security

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用 bcrypt 生成密码哈希
func HashPassword(plain string) (string, error) {
	if plain == "" {
		// 空密码不允许
		return "", errors.New("empty password")
	}
	bs, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

// CheckPassword 校验明文密码与哈希是否匹配
func CheckPassword(hash, plain string) bool {
	if hash == "" || plain == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
