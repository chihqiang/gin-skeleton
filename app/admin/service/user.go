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
	"wangzhiqiang/skeleton/pkg/helper"
)

type UserService struct {
}

// List 获取用户分页列表
func (u UserService) List(ctx context.Context, req *types.UserListReq) (*database.PageResponse[models.SysUser], error) {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return nil, err
	}
	db := apps.DB
	db.Model(models.SysUser{})
	return database.Paginate[models.SysUser](db, req.PageRequest)
}

// Get 获取单个用户
func (u UserService) Get(ctx context.Context, uid uint) (*models.SysUser, error) {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return nil, err
	}
	var user models.SysUser
	if err := apps.DB.Preload("Roles").First(&user, uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Create 创建用户并同步角色到 Casbin
func (u UserService) Create(ctx context.Context, req *types.UserReq) error {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return err
	}
	user := models.SysUser{
		Email:    req.Email,
		Name:     req.Name,
		Phone:    req.Phone,
		Password: req.Password,
	}
	var roles []*models.SysRole
	if len(req.RoleIds) > 0 {
		if err := apps.DB.Where("id IN ?", req.RoleIds).Find(&roles).Error; err != nil {
			return err
		}
	}
	user.Roles = roles
	err = apps.DB.Create(&user).Error
	if err != nil {
		return err
	}
	// 如果分配了角色，同步到 Casbin
	if user.Roles != nil && len(user.Roles) > 0 {
		if err := casbinx.SyncUserRoles(apps.Enforcer, &user); err != nil {
			return err
		}
	}
	return nil
}

// Edit 更新用户信息（可更新角色关联）
func (u UserService) Edit(ctx context.Context, req *types.UserReq) error {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return err
	}
	db := apps.DB
	// 查询原用户
	var updatedUser models.SysUser
	if err := db.Preload("Roles").First(&updatedUser, req.ID).Error; err != nil {
		return err
	}
	updatedUser.Email = req.Email
	updatedUser.Name = req.Name
	updatedUser.Phone = req.Phone
	updatedUser.Password = req.Password
	var roles []*models.SysRole
	if len(req.RoleIds) > 0 {
		if err := apps.DB.Where("id IN ?", req.RoleIds).Find(&roles).Error; err != nil {
			return err
		}
	}
	// 用事务保证更新
	if err := apps.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&updatedUser).Error; err != nil {
			return err
		}
		if roles != nil {
			if err := tx.Model(&updatedUser).Association("Roles").Replace(roles); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	// 事务提交后同步 Casbin
	if roles != nil && len(roles) > 0 {
		if err := casbinx.SyncUserRoles(apps.Enforcer, &updatedUser); err != nil {
			return err
		}
	}
	return nil
}

// Delete 删除用户（安全）
func (u UserService) Delete(ctx context.Context, req *types.IDReq) error {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return err
	}
	if req.ID == apps.Config.System.SuperAdminUID {
		return fmt.Errorf("超级管理员无法删除")
	}
	db := apps.DB
	return db.Transaction(func(tx *gorm.DB) error {
		var user models.SysUser
		if err := tx.Preload("Roles").First(&user, req.ID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil
			}
			return err
		}
		// 移除角色关联
		if err := tx.Model(&user).Association("Roles").Clear(); err != nil {
			return err
		}
		// 删除用户
		if err := tx.Delete(&user).Error; err != nil {
			return err
		}
		// 同步 Casbin：移除该用户所有策略
		if err := casbinx.SyncUserRoles(apps.Enforcer, &user); err != nil {
			return err
		}
		return nil
	})
}

// GetUserMenus 获取用户拥有的菜单（树形结构）
func (u UserService) GetUserMenus(ctx context.Context, uid uint) ([]*models.SysMenu, error) {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return nil, err
	}
	db := apps.DB

	// 超级管理员直接返回所有菜单
	if uid == apps.Config.System.SuperAdminUID {
		var allMenus []*models.SysMenu
		if err := db.Find(&allMenus).Error; err != nil {
			return nil, err
		}
		return helper.BuildTree[*models.SysMenu](allMenus, 0), nil
	}

	// 普通用户：根据角色获取菜单
	var user models.SysUser
	if err := db.Preload("Roles.Menus").First(&user, uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	// 收集用户所有角色关联的菜单
	menuMap := make(map[uint]*models.SysMenu)
	for _, role := range user.Roles {
		for _, menu := range role.Menus {
			menuMap[menu.ID] = menu
		}
	}

	// 转换为切片
	var menus []*models.SysMenu
	for _, menu := range menuMap {
		menus = append(menus, menu)
	}
	// 构建树形结构
	return helper.BuildTree[*models.SysMenu](menus, 0), nil
}
func (u UserService) GetMenus(ctx context.Context, uid uint) ([]*models.SysMenu, error) {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return nil, err
	}
	db := apps.DB
	// 查询所有菜单
	var allMenus []*models.SysMenu
	if err := db.Find(&allMenus).Error; err != nil {
		return nil, err
	}
	selectedMap := make(map[uint]struct{})
	// 超级管理员直接选中所有菜单
	if uid == apps.Config.System.SuperAdminUID {
		for _, menu := range allMenus {
			selectedMap[menu.ID] = struct{}{}
		}
	} else {
		// 普通用户：根据角色获取菜单
		var user models.SysUser
		if err := db.Preload("Roles.Menus").First(&user, uid).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil
			}
			return nil, err
		}
		for _, role := range user.Roles {
			for _, menu := range role.Menus {
				selectedMap[menu.ID] = struct{}{}
			}
		}
	}
	// 转换 map 为切片
	selectedIds := make([]uint, 0, len(selectedMap))
	for id := range selectedMap {
		selectedIds = append(selectedIds, id)
	}
	// 构建树并设置选中状态
	return helper.BuildTreeWithSelected[*models.SysMenu](allMenus, 0, selectedIds), nil
}
