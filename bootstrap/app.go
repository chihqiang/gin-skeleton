package bootstrap

import (
	"wangzhiqiang/skeleton/config"
	"wangzhiqiang/skeleton/pkg/app"
	"wangzhiqiang/skeleton/pkg/httpx"
)

// App 创建并返回一个 app.App 实例
func App(cfg *config.Config) *app.App {
	a := app.New(cfg)
	a.AddPopulate(convertAny[httpx.IController](httpx.GetControllers())...)
	return a
}

// ConvertAny 将特定类型的切片转换为 []any 类型的切片
// 泛型函数，T 可以是任何类型
func convertAny[T any](s []T) []any {
	anyS := make([]any, len(s))
	for i, c := range s {
		anyS[i] = any(c)
	}
	return anyS
}
