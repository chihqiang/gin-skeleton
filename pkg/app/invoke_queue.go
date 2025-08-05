package app

import (
	"context"
	"time"
	"wangzhiqiang/skeleton/config"
	"wangzhiqiang/skeleton/pkg/logger"
	"wangzhiqiang/skeleton/pkg/queue"

	"go.uber.org/fx"
)

type InvokeQueue struct {
	fx.In
	Lc     fx.Lifecycle
	Logger logger.ILogger
	Config *config.Config
	Queue  queue.IQueue
}

func NewInvokeQueue(app InvokeQueue) {
	app.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			startInterval := time.Second
			app.Logger.Infof("[InvokeQueue] Starting queue... interval: %v\n", startInterval)
			go app.Queue.Start(appContext, startInterval)
			app.Logger.Infof("[InvokeQueue] Queue started successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			app.Logger.Infof("[InvokeQueue] Stopping queue...")
			app.Queue.Stop()
			app.Logger.Infof("[InvokeQueue] Queue stopped successfully")
			return nil
		},
	})
}
