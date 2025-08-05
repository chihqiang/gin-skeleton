package admin

import (
	"context"
	"wangzhiqiang/skeleton/app/admin/mock"
	"wangzhiqiang/skeleton/app/admin/models"
	"wangzhiqiang/skeleton/pkg/app"
)

type adminAppInit struct {
}

func (a *adminAppInit) Init(ctx context.Context) {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return
	}
	err = apps.DB.Migrator().AutoMigrate(
		models.SysUser{},
		models.SysMenu{},
		models.SysRole{},
		models.SysAccessLog{},
	)
	if err != nil {
		return
	}
	_ = mock.Save(apps.Enforcer, apps.DB)
}

func init() {
	app.RegisterAppInit(&adminAppInit{})
}
