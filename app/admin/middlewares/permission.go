package middlewares

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"wangzhiqiang/skeleton/pkg/app"
	"wangzhiqiang/skeleton/pkg/casbinx"
	"wangzhiqiang/skeleton/pkg/httpx"
)

func CheckPermission(ctx context.Context) (gin.HandlerFunc, error) {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return nil, err
	}
	enforcer := apps.Enforcer
	return func(c *gin.Context) {
		u, err := GetClaims(c)
		if err != nil {
			httpx.ApiNoAuth(c, err)
			c.Abort()
			return
		}
		allow, _ := casbinx.CheckPermission(enforcer, c.Request, u.UID)
		// 超管直接放行，或者普通用户有权限
		if !allow && u.UID != apps.Config.System.SuperAdminUID {
			httpx.ApiNoForbidden(c, fmt.Errorf("权限不足"))
			c.Abort()
			return
		}
		c.Next()
	}, nil
}
