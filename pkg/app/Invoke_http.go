package app

import (
	"context"
	"time"
	"wangzhiqiang/skeleton/config"
	"wangzhiqiang/skeleton/pkg/httpx"
	"wangzhiqiang/skeleton/pkg/logger"

	"go.uber.org/fx"
)

// ShutdownTimeout 全局变量定义
var (
	ShutdownTimeout = 5 * time.Second // 服务器关闭超时时间
)

// InvokeHTTP 结构体用于依赖注入
// 包含应用生命周期、配置和路由器
type InvokeHTTP struct {
	fx.In                 // 标记为注入类型
	Lc     fx.Lifecycle   // 应用生命周期
	Logger logger.ILogger // 日志记录器
	Config *config.Config // 配置
	Http   *httpx.HTTP
}

// NewInvokeHTTP 设置 HTTP 服务器的生命周期钩子
// 处理服务器的启动和关闭
func NewInvokeHTTP(app InvokeHTTP) {
	app.Lc.Append(fx.Hook{
		// OnStop 钩子：在应用停止时执行
		OnStop: func(ctx context.Context) error {
			app.Logger.Info("[InvokeHTTP] Shutting down server...")
			shutdownCtx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
			defer cancel()
			if err := app.Http.Stop(shutdownCtx); err != nil {
				app.Logger.Errorf("[InvokeHTTP] Error shutting down HTTP server: %v", err)
				return err
			}
			app.Logger.Info("[InvokeHTTP] shut down gracefully")
			return nil
		},
		// OnStart 钩子：在应用启动时执行
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := app.Http.Start(appContext); err != nil {
					app.Logger.Fatalf("[InvokeHTTP] server failed to start %v", err)
				}
			}()
			app.Logger.Info("[InvokeHTTP] start server successfully")
			return nil
		},
	})
}
