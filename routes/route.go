package routes

import (
	_ "wangzhiqiang/skeleton/app/admin"
	"wangzhiqiang/skeleton/app/apis"
	"wangzhiqiang/skeleton/pkg/httpx"
)

func init() {
	httpx.RegisterRoute(&apis.Index{})
	httpx.RegisterRoute(&apis.Task{})
}
