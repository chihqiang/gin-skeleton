package app

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
	"wangzhiqiang/skeleton/config"
	"wangzhiqiang/skeleton/pkg/logger"

	"go.uber.org/fx"
)

// 全局变量定义
var (
	httpServer      *http.Server      // HTTP 服务器实例
	httpServerLock  sync.Mutex        // 保护 httpServer 的互斥锁
	ShutdownTimeout = 5 * time.Second // 服务器关闭超时时间
)

// InvokeHTTP 结构体用于依赖注入
// 包含应用生命周期、配置和路由器
type InvokeHTTP struct {
	fx.In                  // 标记为注入类型
	Lc      fx.Lifecycle   // 应用生命周期
	Logger  logger.ILogger // 日志记录器
	Config  *config.Config // 配置
	Handler http.Handler   // HTTP 处理器
}

// NewInvokeHTTP 设置 HTTP 服务器的生命周期钩子
// 处理服务器的启动和关闭
func NewInvokeHTTP(app InvokeHTTP) {
	app.Lc.Append(fx.Hook{
		// OnStop 钩子：在应用停止时执行
		OnStop: func(ctx context.Context) error {
			app.Logger.Info("[InvokeHTTP] Shutting down server...")
			// 获取当前 HTTP 服务器实例
			httpServerLock.Lock()
			srv := httpServer
			httpServerLock.Unlock()
			if srv == nil {
				app.Logger.Warn("[InvokeHTTP] server instance is nil, skipping shutdown")
				return nil
			}
			// 创建带超时的上下文用于优雅关闭
			shutdownCtx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
			defer cancel()

			// 优雅关闭服务器
			if err := srv.Shutdown(shutdownCtx); err != nil {
				app.Logger.Errorf("[InvokeHTTP] Error shutting down HTTP server: %v", err)
				return err
			}
			app.Logger.Info("[InvokeHTTP] shut down gracefully")
			return nil
		},
		// OnStart 钩子：在应用启动时执行
		OnStart: func(ctx context.Context) error {
			// 创建 HTTP 服务器实例
			srv := &http.Server{
				Addr:         fmt.Sprintf("0.0.0.0:%d", app.Config.Server.Port),           // 服务器地址
				Handler:      app.Handler,                                                 // 处理器
				ReadTimeout:  time.Duration(app.Config.Server.ReadTimeout) * time.Second,  // 读取超时
				WriteTimeout: time.Duration(app.Config.Server.WriteTimeout) * time.Second, // 写入超时
				IdleTimeout:  time.Duration(app.Config.Server.IdleTimeout) * time.Second,  // 空闲超时
			}
			fmt.Println("[InvokeHTTP] server is listening on", srv.Addr)
			// 保存到全局变量
			httpServerLock.Lock()
			httpServer = srv
			httpServerLock.Unlock()
			// 在 goroutine 中启动服务器
			go func() {
				app.Logger.Infof("[InvokeHTTP] start server listening at %s", srv.Addr)
				if err := srv.ListenAndServe(); err != nil {
					fmt.Println("[InvokeHTTP] server failed to start :", err)
					app.Logger.Fatalf("[InvokeHTTP] server failed to start %v", err)
				}
			}()
			app.Logger.Info("[InvokeHTTP] start server successfully")
			return nil
		},
	})
}
