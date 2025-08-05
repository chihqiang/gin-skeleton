package service

import (
	"context"
	"errors"
	"fmt"
	"wangzhiqiang/skeleton/app/admin/models"
	"wangzhiqiang/skeleton/app/admin/types"
	"wangzhiqiang/skeleton/pkg/app"
	"wangzhiqiang/skeleton/pkg/database"
	"wangzhiqiang/skeleton/pkg/helper"

	"gorm.io/gorm"
)

type MenuService struct {
}

// List 分页获取菜单
func (s *MenuService) List(ctx context.Context, req *types.MenuListReq) (*database.PageResponse[models.SysMenu], error) {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return nil, err
	}
	db := apps.DB
	db.Model(models.SysMenu{})
	return database.Paginate[models.SysMenu](db, req.PageRequest)
}

// Tree 获取菜单树
func (s *MenuService) Tree(ctx context.Context) ([]*models.SysMenu, error) {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return nil, err
	}
	db := apps.DB
	var menus []*models.SysMenu
	if err := db.Find(&menus).Error; err != nil {
		return nil, err
	}
	return helper.BuildTree[*models.SysMenu](menus, 0), nil
}

// Get 获取单个菜单
func (s *MenuService) Get(ctx context.Context, id uint) (*models.SysMenu, error) {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return nil, err
	}
	db := apps.DB
	var menu models.SysMenu
	if err := db.First(&menu, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &menu, nil
}

func (s *MenuService) Create(ctx context.Context, req *types.MenuReq) error {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return err
	}
	db := apps.DB
	var menuType models.MenuType
	switch req.Type {
	case "menu":
		menuType = models.MenuTypeMenu // 假设 MenuTypeMenu 是常量值
	case "button":
		menuType = models.MenuTypeButton // 假设 MenuTypeButton 是常量值
	default:
		return fmt.Errorf("无效的菜单类型: %s", req.Type)
	}
	// 构建 SysMenu
	menu := models.SysMenu{
		ParentID:  req.ParentID,
		Name:      req.Name,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		Sort:      req.Sort,
		Hidden:    req.Hidden,
		Type:      menuType,
		Method:    req.Method,
	}
	// 解析角色
	var roles []models.SysRole
	if len(req.RoleIds) > 0 {
		if err := db.Where("id IN ?", req.RoleIds).Find(&roles).Error; err != nil {
			return err
		}
	}
	// 创建菜单 + 绑定角色
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&menu).Error; err != nil {
			return err
		}
		if len(roles) > 0 {
			if err := tx.Model(&menu).Association("Roles").Replace(roles); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *MenuService) Edit(ctx context.Context, req *types.MenuReq) error {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return err
	}
	var menuType models.MenuType
	switch req.Type {
	case "menu":
		menuType = models.MenuTypeMenu // 假设 MenuTypeMenu 是常量值
	case "button":
		menuType = models.MenuTypeButton // 假设 MenuTypeButton 是常量值
	default:
		return fmt.Errorf("无效的菜单类型: %s", req.Type)
	}
	db := apps.DB
	// 查询原菜单
	var menu models.SysMenu
	if err := db.Preload("Roles").First(&menu, req.ID).Error; err != nil {
		return err
	}
	menu.ParentID = req.ParentID
	menu.Name = req.Name
	menu.Path = req.Path
	menu.Component = req.Component
	menu.Icon = req.Icon
	menu.Sort = req.Sort
	menu.Hidden = req.Hidden
	menu.Method = req.Method
	menu.Type = menuType
	// 解析角色
	var roles []models.SysRole
	updateRoles := false
	if len(req.RoleIds) > 0 {
		updateRoles = true
		if err := db.Where("id IN ?", req.RoleIds).Find(&roles).Error; err != nil {
			return err
		}
	}

	// 用事务更新菜单和角色关联
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&menu).Error; err != nil {
			return err
		}
		if updateRoles {
			if err := tx.Model(&menu).Association("Roles").Replace(roles); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

// Delete 安全删除菜单
func (s *MenuService) Delete(ctx context.Context, id uint) error {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return err
	}
	db := apps.DB

	var menu models.SysMenu
	if err := db.Preload("Roles").First(&menu, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // 菜单不存在，可认为删除成功
		}
		return err
	}
	// 检查是否关联角色
	if len(menu.Roles) > 0 {
		return errors.New("该菜单已被角色关联，无法删除")
	}
	// 可选：检查是否有子菜单，防止破坏树结构
	var childCount int64
	db.Model(&models.SysMenu{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		return errors.New("该菜单存在子菜单，无法删除")
	}
	return db.Delete(&models.SysMenu{}, id).Error
}
