package apis

import (
	"context"
	"github.com/gin-gonic/gin"
	"wangzhiqiang/skeleton/app/tasks"
	"wangzhiqiang/skeleton/pkg/app"
	"wangzhiqiang/skeleton/pkg/httpx"
)

type Task struct {
}

func (c *Task) Routes(ctx context.Context, g *gin.Engine) error {
	g.GET("/test/queue", func(context *gin.Context) {
		apps, err := app.GetApps(ctx)
		if err != nil {
			httpx.ApiError(context, err)
			return
		}
		if err := apps.Queue.Push(&tasks.EmailTask{
			To:      "test@example.com",
			Subject: "测试邮件",
			Body:    "这是邮件内容",
		}, 0); err != nil {
			httpx.ApiError(context, err)
			return
		}
		httpx.ApiSuccess(context, "邮件已发送到队列")
	})
	return nil
}
