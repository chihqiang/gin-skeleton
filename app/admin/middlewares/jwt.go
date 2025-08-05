package middlewares

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"wangzhiqiang/skeleton/pkg/httpx"
	"wangzhiqiang/skeleton/pkg/jwts"
)

const CtxKeyUserClaims = "claims"

// GetClaims 从 gin.Context 获取 JWT 用户信息
func GetClaims(c *gin.Context) (*jwts.Claims, error) {
	user, exists := c.Get(CtxKeyUserClaims)
	if !exists {

		return nil, fmt.Errorf("未登录")
	}
	claims, ok := user.(*jwts.Claims)
	if !ok {
		return nil, fmt.Errorf("用户信息解析失败,请重新登录")
	}
	return claims, nil
}

// JWTAuth 返回 JWT 认证中间件
func JWTAuth(jwt *jwts.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			httpx.ApiNoAuth(c, errors.New("缺少 token"))
			c.Abort()
			return
		}

		claims, err := jwt.Parse(token)
		if err != nil {
			if errors.Is(err, jwts.TokenExpired) {
				httpx.ApiNoAuth(c, errors.New("登录已过期，请重新登录"))
				c.Abort()
				return
			}
			httpx.ApiNoAuth(c, fmt.Errorf("token 无效: %w", err))
			c.Abort()
			return
		}

		// 写入 context，供后续中间件使用
		c.Set(CtxKeyUserClaims, claims)
		c.Next()
	}
}
