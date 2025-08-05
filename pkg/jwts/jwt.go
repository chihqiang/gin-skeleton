package jwts

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenType
// -------------------- 常量与错误 --------------------
type TokenType string

// 编译时检查是否实现接口
var _ JWTService = (*JWT)(nil)

const (
	AccessTokenType  TokenType = "access_token"
	RefreshTokenType TokenType = "refresh_token"
)

var (
	TokenValid            = errors.New("未知错误")
	TokenExpired          = errors.New("token已过期")
	TokenNotValidYet      = errors.New("token尚未激活")
	TokenMalformed        = errors.New("这不是一个token")
	TokenSignatureInvalid = errors.New("无效签名")
	TokenInvalid          = errors.New("无法处理此token")
	ErrNotRefreshToken    = errors.New("不是 refresh token，无法续签")
)

// IUser
// -------------------- 用户接口 --------------------
type IUser interface {
	GetID() uint
	GetName() string
}

// Claims
// -------------------- Claims --------------------
type Claims struct {
	UID       uint      `json:"uid"`
	Name      string    `json:"name"`
	TokenType TokenType `json:"token_type"` // access_token / refresh_token
	jwt.RegisteredClaims
}

// Config
// -------------------- 配置 --------------------
type Config struct {
	Secret            string   `yaml:"secret" json:"secret,omitempty"`
	Expiration        int      `yaml:"expiration" json:"expiration,omitempty"`                 // Access Token 过期秒数，0 表示永不过期
	RefreshExpiration int      `yaml:"refresh_expiration" json:"refresh_expiration,omitempty"` // Refresh Token 过期秒数，0 表示永不过期
	Issuer            string   `yaml:"issuer" json:"issuer,omitempty"`
	Aud               []string `yaml:"aud" json:"aud,omitempty"`
}

// JWTService
// -------------------- JWT 接口 --------------------
type JWTService interface {
	BuildAccessToken(user IUser) (string, error)
	BuildRefreshToken(user IUser) (string, error)
	Parse(tokenStr string) (*Claims, error)
	Refresh(rToken string) (accessToken string, refreshToken string, err error)
	GetAccessExpiresIn() int
	GetRefreshExpiresIn() int
	ExpiresIn(tokenStr string) (int, error)
}

// -------------------- JWT 实现 --------------------
type JWT struct {
	config *Config
}

func NewJWT(cfg *Config) *JWT {
	if cfg.Issuer == "" {
		cfg.Issuer = "skeleton"
	}
	if len(cfg.Aud) == 0 {
		cfg.Aud = []string{"skeleton"}
	}
	return &JWT{config: cfg}
}

// 内部通用生成 token
func (j *JWT) buildToken(uid uint, name string, tokenType TokenType, expSec int) (string, error) {
	now := time.Now()
	claims := Claims{
		UID:       uid,
		Name:      name,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        fmt.Sprintf("%d", uid),
			Subject:   fmt.Sprintf("%d", uid),
			Issuer:    j.config.Issuer,
			Audience:  jwt.ClaimStrings(j.config.Aud),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	if expSec > 0 {
		claims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Duration(expSec) * time.Second))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.Secret))
}

// BuildAccessToken
// -------------------- 对外方法 --------------------
func (j *JWT) BuildAccessToken(user IUser) (string, error) {
	return j.buildToken(user.GetID(), user.GetName(), AccessTokenType, j.config.Expiration)
}

func (j *JWT) BuildRefreshToken(user IUser) (string, error) {
	return j.buildToken(user.GetID(), user.GetName(), RefreshTokenType, j.config.RefreshExpiration)
}

func (j *JWT) Parse(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(j.config.Secret), nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, TokenExpired
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, TokenMalformed
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, TokenSignatureInvalid
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, TokenNotValidYet
		default:
			return nil, TokenInvalid
		}
	}
	if !token.Valid {
		return nil, TokenValid
	}
	return claims, nil
}

// Refresh 用 refresh token 续签新的 access + refresh
func (j *JWT) Refresh(refreshToken string) (string, string, error) {
	claims, err := j.Parse(refreshToken)
	if err != nil {
		return "", "", err
	}
	if claims.TokenType != RefreshTokenType {
		return "", "", ErrNotRefreshToken
	}

	newAccess, err := j.buildToken(claims.UID, claims.Name, AccessTokenType, j.config.Expiration)
	if err != nil {
		return "", "", err
	}
	newRefresh, err := j.buildToken(claims.UID, claims.Name, RefreshTokenType, j.config.RefreshExpiration)
	if err != nil {
		return "", "", err
	}
	return newAccess, newRefresh, nil
}

// GetAccessExpiresIn 返回 Access Token 配置过期秒数
func (j *JWT) GetAccessExpiresIn() int {
	return j.config.Expiration
}

// GetRefreshExpiresIn 返回 Refresh Token 配置过期秒数
func (j *JWT) GetRefreshExpiresIn() int {
	return j.config.RefreshExpiration
}

// ExpiresIn 动态计算 token 剩余秒数
func (j *JWT) ExpiresIn(tokenStr string) (int, error) {
	claims, err := j.Parse(tokenStr)
	if err != nil {
		return 0, err
	}
	if claims.ExpiresAt == nil {
		return -1, nil
	}
	sec := int(time.Until(claims.ExpiresAt.Time).Seconds())
	if sec < 0 {
		sec = 0
	}
	return sec, nil
}
