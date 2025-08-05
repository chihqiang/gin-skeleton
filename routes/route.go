package routes

import (
	"wangzhiqiang/skeleton/app/admin"
	"wangzhiqiang/skeleton/app/apis"
	"wangzhiqiang/skeleton/pkg/httpx"
)

func init() {
	httpx.RegisterRoute(&apis.Index{})
	httpx.RegisterRoute(&apis.Task{})
	admin.LoadRoute()
}
