package admin

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"wangzhiqiang/skeleton/app/admin/apis"
	"wangzhiqiang/skeleton/app/admin/middlewares"
	"wangzhiqiang/skeleton/pkg/httpx"
	"wangzhiqiang/skeleton/pkg/httpx/mws"
)

type adminRoute struct {
}

func (a *adminRoute) Routes(ctx context.Context, g *gin.Engine) error {
	g.NoRoute(func(c *gin.Context) {
		httpx.ApiErrWithCode(c, fmt.Errorf("not Found"), http.StatusNotFound)
	})
	jwtAuth, err := middlewares.JWTAuth(ctx)
	if err != nil {
		return err
	}
	permission, err := middlewares.CheckPermission(ctx)
	if err != nil {
		return err
	}
	g.Use(mws.Core())
	adminGroup := g.Group("/api/admin")
	{
		// 登录不需要认证
		adminGroup.POST("/login", apis.NewAuth(ctx).Login)
		// 刷新token
		adminGroup.POST("/refresh", apis.NewAuth(ctx).Refresh)

		// 用户和角色操作需要认证和权限中间件
		adminGroup.Use(jwtAuth)
		adminGroup.Use(permission)
		// 用户管理
		user := apis.NewUser(ctx)
		userGroup := adminGroup.Group("/user")
		{
			userGroup.GET("", user.List)             // 查询用户列表
			userGroup.POST("/create", user.Create)   // 创建用户
			userGroup.GET("/view", user.View)        // 查看用户
			userGroup.PUT("/edit", user.Edit)        // 编辑用户
			userGroup.DELETE("/delete", user.Delete) // 删除用户
		}

		// 角色管理
		role := apis.NewRole(ctx)
		roleGroup := adminGroup.Group("/role")
		{
			roleGroup.GET("", role.List)               // 查询角色列表
			roleGroup.POST("/create", role.Create)     // 创建角色
			roleGroup.GET("/view", role.View)          // 查看角色
			roleGroup.PUT("/edit", role.Edit)          // 编辑角色
			roleGroup.DELETE("/delete", role.Delete)   // 删除角色
			roleGroup.GET("/auth/list", role.AuthList) // 角色权限列表
			roleGroup.POST("/auth", role.Auth)         // 授权权限
		}

		// 菜单管理
		menu := apis.NewMenu(ctx)
		menuGroup := adminGroup.Group("/menu")
		{
			menuGroup.GET("", menu.List)             // 查询菜单列表
			menuGroup.POST("/create", menu.Create)   // 创建菜单
			menuGroup.GET("/view", menu.View)        // 查看菜单
			menuGroup.PUT("/edit", menu.Edit)        // 编辑菜单
			menuGroup.DELETE("/delete", menu.Delete) // 删除菜单
		}
	}
	return nil
}

func LoadRoute() {
	httpx.RegisterRoute(&adminRoute{})
}
