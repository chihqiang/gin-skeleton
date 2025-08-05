package admin

import (
	"context"
	"wangzhiqiang/skeleton/app/admin/mock"
	"wangzhiqiang/skeleton/app/admin/models"
	appModels "wangzhiqiang/skeleton/app/models"
	"wangzhiqiang/skeleton/pkg/app"
	"wangzhiqiang/skeleton/pkg/httpx"
)

type adminAppInit struct {
}

func (a *adminAppInit) Init(ctx context.Context) {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return
	}
	_ = apps.DB.Migrator().AutoMigrate(
		models.SysUser{},
		models.SysMenu{},
		models.SysRole{},
		appModels.SysAccessLog{},
	)
	_ = mock.Save(apps.Enforcer, apps.DB)
}

func init() {
	app.RegisterAppInit(&adminAppInit{})

	httpx.RegisterRoute(&Route{})
}
