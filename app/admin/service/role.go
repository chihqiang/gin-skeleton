package service

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"wangzhiqiang/skeleton/app/admin/models"
	"wangzhiqiang/skeleton/app/admin/types"
	"wangzhiqiang/skeleton/pkg/app"
	"wangzhiqiang/skeleton/pkg/casbinx"
	"wangzhiqiang/skeleton/pkg/database"
)

type RoleService struct {
}

// List 获取角色分页列表
func (s *RoleService) List(ctx context.Context, req *types.RoleListReq) (*database.PageResponse[models.SysRole], error) {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return nil, err
	}
	db := apps.DB
	db.Model(models.SysRole{})
	return database.Paginate[models.SysRole](db, req.PageRequest)
}

// Create 创建角色
func (s *RoleService) Create(ctx context.Context, req *types.RoleReq) error {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return err
	}
	db := apps.DB
	role := models.SysRole{
		Name:   req.Name,
		Code:   req.Code,
		Remark: req.Remark,
	}

	// 查询菜单对象并关联（即便为空也可）
	var menus []*models.SysMenu
	if err := db.Where("id IN ?", req.MenuIds).Find(&menus).Error; err != nil {
		return err
	}
	role.Menus = menus

	if err := db.Create(&role).Error; err != nil {
		return err
	}

	// 创建后同步 Casbin
	return casbinx.SyncRole(apps.Enforcer, &role)
}

// View 查看角色详情
func (s *RoleService) View(ctx context.Context, id uint) (*models.SysRole, error) {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return nil, err
	}
	var role models.SysRole
	if err := apps.DB.Preload("Menus").First(&role, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("角色不存在")
		}
		return nil, err
	}
	return &role, nil
}

// Edit 更新角色信息及菜单关联
func (s *RoleService) Edit(ctx context.Context, req *types.RoleReq) error {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return err
	}
	db := apps.DB

	var role models.SysRole
	if err := db.Preload("Menus").First(&role, req.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("角色不存在")
		}
		return err
	}

	role.Name = req.Name
	role.Code = req.Code
	role.Remark = req.Remark

	if err := db.Save(&role).Error; err != nil {
		return err
	}

	// 查询菜单对象并替换关联（即便 MenuIds 为空也会清空关联）
	var menus []*models.SysMenu
	if err := db.Where("id IN ?", req.MenuIds).Find(&menus).Error; err != nil {
		return err
	}
	if err := db.Model(&role).Association("Menus").Replace(menus); err != nil {
		return err
	}

	// 更新后同步 Casbin
	return casbinx.SyncRole(apps.Enforcer, &role)
}

// Delete 删除角色及菜单关联
func (s *RoleService) Delete(ctx context.Context, id uint) error {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return err
	}
	db := apps.DB

	var role models.SysRole
	if err := db.First(&role, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("角色不存在")
		}
		return err
	}

	// 清空菜单关联
	if err := db.Model(&role).Association("Menus").Clear(); err != nil {
		return err
	}

	if err := db.Delete(&role).Error; err != nil {
		return err
	}

	// 删除后同步 Casbin
	return casbinx.SyncRole(apps.Enforcer, &role)
}

// Auth 授权角色菜单
func (s *RoleService) Auth(ctx context.Context, req *types.RoleAuthReq) error {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return err
	}
	db := apps.DB

	var role models.SysRole
	if err := db.First(&role, req.RoleID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("角色不存在")
		}
		return err
	}

	// 查询菜单对象，替换关联，即便 MenuIDs 为空也会清空
	var menus []*models.SysMenu
	if err := db.Where("id IN ?", req.MenuIDs).Find(&menus).Error; err != nil {
		return err
	}
	if err := db.Model(&role).Association("Menus").Replace(menus); err != nil {
		return err
	}

	// 同步到 Casbin
	return casbinx.SyncRole(apps.Enforcer, &role)
}
