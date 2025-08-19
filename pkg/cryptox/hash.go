package cryptox

import (
	"golang.org/x/crypto/bcrypt"
)

// HashMake 生成密码哈希
func HashMake(raw string) string {
	pwd := []byte(raw)
	hash, _ := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	return string(hash)
}

// HashVerify 验证密码是否匹配哈希
func HashVerify(raw, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw)) == nil
}
