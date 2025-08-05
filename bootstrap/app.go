package bootstrap

import (
	"wangzhiqiang/skeleton/config"
	"wangzhiqiang/skeleton/pkg/app"
)

// App 创建并返回一个 app.App 实例
func App(cfg *config.Config) *app.App {
	a := app.New(cfg)
	return a
}
