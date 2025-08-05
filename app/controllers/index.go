package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"wangzhiqiang/skeleton/app/tasks"
	"wangzhiqiang/skeleton/config"
	"wangzhiqiang/skeleton/pkg/httpx"
	"wangzhiqiang/skeleton/pkg/queue"
)

type Index struct {
	fx.In
	Config *config.Config
	Queue  queue.IQueue
}

func (index *Index) Routes(g gin.IRoutes) {
	g.GET("/", func(context *gin.Context) {
		err := index.Queue.Push(&tasks.EmailTask{
			To:      "test@example.com",
			Subject: "测试邮件",
			Body:    "这是邮件内容",
		}, 0)
		if err != nil {
			httpx.Error(context, err)
			return
		}
		httpx.Success(context, index.Config)
	})
}
