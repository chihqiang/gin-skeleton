package app

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"sync"
	"wangzhiqiang/skeleton/config"
	"wangzhiqiang/skeleton/pkg/logger"
)

var (
	ContextAppKey = "app"
	ctxApp        Apps
	ctxAppLock    sync.Mutex
	appContext    context.Context
)

type Apps struct {
	fx.In
	Lc     fx.Lifecycle
	Logger logger.ILogger
	Config *config.Config
	Redis  *redis.Client
}

func NewApps(app Apps) {
	app.Lc.Append(fx.Hook{
		// OnStop 钩子：在应用停止时执行
		OnStop: func(ctx context.Context) error {
			return nil
		},
		// OnStart 钩子：在应用启动时执行
		OnStart: func(ctx context.Context) error {
			// 保存到全局变量
			ctxAppLock.Lock()
			ctxApp = app
			appContext = context.WithValue(ctx, ContextAppKey, ctxApp)
			ctxAppLock.Unlock()
			return nil
		},
	})
}
