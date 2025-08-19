package admin

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"wangzhiqiang/skeleton/app/admin/apis"
	"wangzhiqiang/skeleton/app/admin/middlewares"
	appMiddlewares "wangzhiqiang/skeleton/app/middlewares"
	"wangzhiqiang/skeleton/pkg/app"
	"wangzhiqiang/skeleton/pkg/httpx"
	"wangzhiqiang/skeleton/pkg/httpx/mws"
)

type Route struct {
}

func (a *Route) Routes(ctx context.Context, g *gin.Engine) error {
	g.NoRoute(func(c *gin.Context) {
		httpx.ApiErrWithCode(c, fmt.Errorf("not Found"), http.StatusNotFound)
	})
	apps, err := app.GetApps(ctx)
	if err != nil {
		return err
	}
	jwtAuth := middlewares.JWTAuth(apps.JWT)
	permission := middlewares.CheckPermission(apps.Enforcer, apps.Config)
	accessLog := appMiddlewares.AccessLog(apps.DB)
	g.Use(mws.Core())
	api := apis.NewApis(ctx)
	adminGroup := g.Group("/api/admin")
	{
		// 登录不需要认证
		adminGroup.POST("/login", api.Auth.Login)
		// 用户和角色操作需要认证和权限中间件
		adminGroup.Use(jwtAuth)
		adminGroup.Use(accessLog)
		adminGroup.Use(permission)

		// 刷新token
		adminGroup.POST("/refresh", api.Auth.Refresh)
		// 用户管理
		userGroup := adminGroup.Group("/user")
		{
			userGroup.GET("", api.User.List)             // 查询用户列表
			userGroup.POST("/create", api.User.Create)   // 创建用户
			userGroup.GET("/view", api.User.View)        // 查看用户
			userGroup.PUT("/edit", api.User.Edit)        // 编辑用户
			userGroup.DELETE("/delete", api.User.Delete) // 删除用户
		}

		// 角色管理
		roleGroup := adminGroup.Group("/role")
		{
			roleGroup.GET("", api.Role.List)               // 查询角色列表
			roleGroup.POST("/create", api.Role.Create)     // 创建角色
			roleGroup.GET("/view", api.Role.View)          // 查看角色
			roleGroup.PUT("/edit", api.Role.Edit)          // 编辑角色
			roleGroup.DELETE("/delete", api.Role.Delete)   // 删除角色
			roleGroup.GET("/auth/list", api.Role.AuthList) // 角色权限列表
			roleGroup.POST("/auth", api.Role.Auth)         // 授权权限
		}

		// 菜单管理
		menuGroup := adminGroup.Group("/menu")
		{
			menuGroup.GET("", api.Menu.List)             // 查询菜单列表
			menuGroup.POST("/create", api.Menu.Create)   // 创建菜单
			menuGroup.GET("/view", api.Menu.View)        // 查看菜单
			menuGroup.PUT("/edit", api.Menu.Edit)        // 编辑菜单
			menuGroup.DELETE("/delete", api.Menu.Delete) // 删除菜单
		}
	}
	return nil
}
