package app

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"sync"
	"wangzhiqiang/skeleton/config"
	"wangzhiqiang/skeleton/pkg/jwts"
	"wangzhiqiang/skeleton/pkg/logger"
	"wangzhiqiang/skeleton/pkg/queue"
)

var (
	_initApps []IAppInit
)

// 定义一个私有类型，避免与其他 package 冲突
type contextKey string

const ContextAppKey contextKey = "appKey"

var (
	ctxApp     Apps
	ctxAppLock sync.Mutex
	appContext context.Context
)

type IAppInit interface {
	Init(ctx context.Context)
}

func RegisterAppInit(initApp IAppInit) {
	_initApps = append(_initApps, initApp)
}

type Apps struct {
	fx.In
	Lc       fx.Lifecycle
	Logger   logger.ILogger
	Config   *config.Config
	Queue    queue.IQueue
	DB       *gorm.DB
	JWT      *jwts.JWT
	Enforcer *casbin.Enforcer
}

func GetApps(ctx context.Context) (Apps, error) {
	if apps, ok := ctx.Value(ContextAppKey).(Apps); ok {
		return apps, nil
	}
	return Apps{}, fmt.Errorf("apps not found in context")
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
			for _, initApp := range _initApps {
				initApp.Init(appContext)
			}
			ctxAppLock.Unlock()
			return nil
		},
	})
}
