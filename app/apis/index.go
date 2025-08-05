package apis

import (
	"context"
	"github.com/gin-gonic/gin"
	"wangzhiqiang/skeleton/pkg/app"
	"wangzhiqiang/skeleton/pkg/httpx"
)

type Index struct {
}

func (index *Index) Routes(ctx context.Context, g *gin.Engine) error {
	g.GET("/", func(context *gin.Context) {
		apps, err := app.GetApps(ctx)
		if err != nil {
			httpx.ApiError(context, err)
			return
		}
		httpx.ApiSuccess(context, apps.Config)
	})
	return nil
}
