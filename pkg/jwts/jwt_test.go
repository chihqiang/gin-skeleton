package jwts

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mock 用户实现 IUser 接口
type mockUser struct {
	id   uint
	name string
}

func (u *mockUser) GetID() uint     { return u.id }
func (u *mockUser) GetName() string { return u.name }

func TestJWT_Basic(t *testing.T) {
	user := &mockUser{id: 123, name: "Alice"}

	jwtCfg := &Config{
		Secret:            "test-secret",
		Expiration:        3600, // 1小时
		RefreshExpiration: 7200, // 2小时
		Issuer:            "test-issuer",
		Aud:               []string{"test-aud"},
	}

	j := NewJWT(jwtCfg)

	// 生成 access token
	accessToken, err := j.BuildAccessToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)

	// 生成 refresh token
	refreshToken, err := j.BuildRefreshToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshToken)

	// 解析 access token
	accessClaims, err := j.Parse(accessToken)
	assert.NoError(t, err)
	assert.Equal(t, AccessTokenType, accessClaims.TokenType)
	assert.Equal(t, user.GetID(), accessClaims.UID)
	assert.Equal(t, user.GetName(), accessClaims.Name)

	// 解析 refresh token
	refreshClaims, err := j.Parse(refreshToken)
	assert.NoError(t, err)
	assert.Equal(t, RefreshTokenType, refreshClaims.TokenType)
	assert.Equal(t, user.GetID(), refreshClaims.UID)
	assert.Equal(t, user.GetName(), refreshClaims.Name)
}

func TestJWT_Refresh(t *testing.T) {
	user := &mockUser{id: 123, name: "Alice"}

	jwtCfg := &Config{
		Secret:            "test-secret",
		Expiration:        1,    // 1秒，快速过期
		RefreshExpiration: 3600, // 1小时
		Issuer:            "test-issuer",
		Aud:               []string{"test-aud"},
	}

	j := NewJWT(jwtCfg)

	// 构造 refresh token
	refreshToken, err := j.BuildRefreshToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshToken)

	// 用 refresh token 续签
	newAccess, newRefresh, err := j.Refresh(refreshToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, newAccess)
	assert.NotEmpty(t, newRefresh)

	// refresh token 续签 access token 类型检查
	newAccessClaims, err := j.Parse(newAccess)
	assert.NoError(t, err)
	assert.Equal(t, AccessTokenType, newAccessClaims.TokenType)
	assert.Equal(t, user.GetID(), newAccessClaims.UID)

	newRefreshClaims, err := j.Parse(newRefresh)
	assert.NoError(t, err)
	assert.Equal(t, RefreshTokenType, newRefreshClaims.TokenType)
	assert.Equal(t, user.GetID(), newRefreshClaims.UID)
}

func TestJWT_ExpiresIn(t *testing.T) {
	user := &mockUser{id: 123, name: "Alice"}

	jwtCfg := &Config{
		Secret:            "test-secret",
		Expiration:        2, // 2秒
		RefreshExpiration: 3, // 3秒
		Issuer:            "test-issuer",
		Aud:               []string{"test-aud"},
	}

	j := NewJWT(jwtCfg)

	accessToken, _ := j.BuildAccessToken(user)
	refreshToken, _ := j.BuildRefreshToken(user)

	// Access token 剩余时间
	sec, err := j.ExpiresIn(accessToken)
	assert.NoError(t, err)
	assert.True(t, sec <= 2 && sec > 0)

	// Refresh token 剩余时间
	sec, err = j.ExpiresIn(refreshToken)
	assert.NoError(t, err)
	assert.True(t, sec <= 3 && sec > 0)

	// 等待 access token 过期
	time.Sleep(3 * time.Second)
	_, err = j.Parse(accessToken)
	assert.Equal(t, TokenExpired, err)
}
