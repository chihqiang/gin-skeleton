package middlewares

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"wangzhiqiang/skeleton/config"
	"wangzhiqiang/skeleton/pkg/casbinx"
	"wangzhiqiang/skeleton/pkg/httpx"
)

func CheckPermission(enforcer *casbin.Enforcer, config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, err := GetClaims(c)
		if err != nil {
			httpx.ApiNoAuth(c, err)
			c.Abort()
			return
		}
		allow, _ := casbinx.CheckPermission(enforcer, c.Request, u.UID)
		// 超管直接放行，或者普通用户有权限
		if !allow && u.UID != config.System.SuperAdminUID {
			httpx.ApiNoForbidden(c, fmt.Errorf("权限不足"))
			c.Abort()
			return
		}
		c.Next()
	}
}
