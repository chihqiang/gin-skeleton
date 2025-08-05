package httpx

import (
	"context"
	"github.com/gin-gonic/gin"
)

var (
	routes []IRoute
)

// IRoute 路由
// 所有控制器都需要实现这个接口来注册路由
type IRoute interface {
	// Routes 注册控制器的路由
	// 参数 g: Gin 路由接口
	Routes(ctx context.Context, g *gin.Engine) error
}

// RegisterRoute 注册一个控制器实例
// 通常在 init 函数或模块初始化中调用
// 所有注册的控制器会保存在 routes 切片中
// 参数 controller: 要注册的控制器实例
func RegisterRoute(route IRoute) {
	routes = append(routes, route)
}
