package routes

import (
	"wangzhiqiang/skeleton/app/controllers"
	"wangzhiqiang/skeleton/pkg/httpx"
)

func init() {
	httpx.RegisterController(&controllers.Index{})
}
