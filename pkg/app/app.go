package app

import (
	"wangzhiqiang/skeleton/config"

	"go.uber.org/fx"
)

type App struct {
	opts []fx.Option
}

// New 创建一个新的 App 实例
func New(cfg *config.Config) *App {
	app := &App{opts: make([]fx.Option, 0)}

	app.AddProvide(
		func() *config.Config { return cfg },
		ProvideLogger,     // 提供日志记录器
		ProvideDatabase,   // 提供数据库
		ProvideEnforcer,   // 提供Casbin
		ProvideHTTPServer, // 提供服务器
		ProvideQueue,      // 提供队列
		ProvideJWT,        // 提供JWT服务
	)
	//if cfg.Server.Mode != "debug" {
	app.AddOpts(fx.NopLogger)
	//}
	app.AddInvoke(NewApps)
	return app
}

// AddPopulate 添加需要填充的目标
func (a *App) AddPopulate(targets ...any) {
	a.AddOpts(fx.Populate(targets...))
}

// AddProvide 添加提供者函数
func (a *App) AddProvide(constructors ...any) {
	a.AddOpts(fx.Provide(constructors...))
}

// AddInvoke 添加需要调用的函数
func (a *App) AddInvoke(invoke any) {
	a.AddOpts(fx.Invoke(invoke))
}

// AddOpts 添加额外的 fx.Option
func (a *App) AddOpts(opts ...fx.Option) {
	a.opts = append(a.opts, opts...)
}

// FX 返回一个 fx.App 实例
func (a *App) FX() *fx.App {
	return fx.New(a.opts...)
}

// StartHTTP 启动 HTTP 服务器
func (a *App) StartHTTP() {
	a.AddInvoke(NewInvokeHTTP)
	app := a.FX()
	app.Run()
}

// StartQueue 启动队列处理器
func (a *App) StartQueue() {
	a.AddInvoke(NewInvokeQueue)
	app := a.FX()
	app.Run()
}
