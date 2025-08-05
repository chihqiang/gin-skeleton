package mock

import (
	"github.com/casbin/casbin/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"net/http"
	"wangzhiqiang/skeleton/app/admin/models"
	"wangzhiqiang/skeleton/pkg/casbinx"
	"wangzhiqiang/skeleton/pkg/cryptox"
)

func Save(e *casbin.Enforcer, db *gorm.DB) error {
	prefix := "/api/admin"

	// ---------------- 创建主菜单 ----------------
	dashboard := models.SysMenu{Name: "首页", Path: prefix + "/dashboard", Method: datatypes.JSONSlice[string]{http.MethodGet}, Sort: 1, Type: models.MenuTypeMenu}
	usersMenu := models.SysMenu{Name: "用户管理", Path: prefix + "/user", Method: datatypes.JSONSlice[string]{http.MethodGet}, Sort: 2, Type: models.MenuTypeMenu}
	rolesMenu := models.SysMenu{Name: "角色管理", Path: prefix + "/role", Method: datatypes.JSONSlice[string]{http.MethodGet}, Sort: 3, Type: models.MenuTypeMenu}
	menusMenu := models.SysMenu{Name: "菜单管理", Path: prefix + "/menu", Method: datatypes.JSONSlice[string]{http.MethodGet}, Sort: 4, Type: models.MenuTypeMenu}

	menus := []*models.SysMenu{&dashboard, &usersMenu, &rolesMenu, &menusMenu}
	for _, m := range menus {
		if err := db.Create(m).Error; err != nil {
			return err
		}
	}

	// ---------------- 用户管理操作 ----------------
	userCreate := models.SysMenu{Name: "创建用户", Path: prefix + "/user/create", Method: datatypes.JSONSlice[string]{http.MethodPost}, ParentID: usersMenu.ID, Type: models.MenuTypeButton}
	userView := models.SysMenu{Name: "查看用户", Path: prefix + "/user/view", Method: datatypes.JSONSlice[string]{http.MethodGet}, ParentID: usersMenu.ID, Type: models.MenuTypeButton}
	userEdit := models.SysMenu{Name: "编辑用户", Path: prefix + "/user/edit", Method: datatypes.JSONSlice[string]{http.MethodPut}, ParentID: usersMenu.ID, Type: models.MenuTypeButton}
	userDelete := models.SysMenu{Name: "删除用户", Path: prefix + "/user/delete", Method: datatypes.JSONSlice[string]{http.MethodDelete}, ParentID: usersMenu.ID, Type: models.MenuTypeButton}

	// ---------------- 角色管理操作 ----------------
	roleCreate := models.SysMenu{Name: "创建角色", Path: prefix + "/role/create", Method: datatypes.JSONSlice[string]{http.MethodPost}, ParentID: rolesMenu.ID, Type: models.MenuTypeButton}
	roleView := models.SysMenu{Name: "查看角色", Path: prefix + "/role/view", Method: datatypes.JSONSlice[string]{http.MethodGet}, ParentID: rolesMenu.ID, Type: models.MenuTypeButton}
	roleEdit := models.SysMenu{Name: "编辑角色", Path: prefix + "/role/edit", Method: datatypes.JSONSlice[string]{http.MethodPut}, ParentID: rolesMenu.ID, Type: models.MenuTypeButton}
	roleDelete := models.SysMenu{Name: "删除角色", Path: prefix + "/role/delete", Method: datatypes.JSONSlice[string]{http.MethodDelete}, ParentID: rolesMenu.ID, Type: models.MenuTypeButton}
	roleAuth := models.SysMenu{Name: "授权权限", Path: prefix + "/role/auth", Method: datatypes.JSONSlice[string]{http.MethodPost}, ParentID: rolesMenu.ID, Type: models.MenuTypeButton}
	roleAuthList := models.SysMenu{Name: "权限列表", Path: prefix + "/role/auth/list", Method: datatypes.JSONSlice[string]{http.MethodGet}, ParentID: rolesMenu.ID, Type: models.MenuTypeButton}

	// ---------------- 菜单管理操作 ----------------
	menuCreate := models.SysMenu{Name: "创建菜单", Path: prefix + "/menu/create", Method: datatypes.JSONSlice[string]{http.MethodPost}, ParentID: menusMenu.ID, Type: models.MenuTypeButton}
	menuView := models.SysMenu{Name: "查看菜单", Path: prefix + "/menu/view", Method: datatypes.JSONSlice[string]{http.MethodGet}, ParentID: menusMenu.ID, Type: models.MenuTypeButton}
	menuEdit := models.SysMenu{Name: "编辑菜单", Path: prefix + "/menu/edit", Method: datatypes.JSONSlice[string]{http.MethodPut}, ParentID: menusMenu.ID, Type: models.MenuTypeButton}
	menuDelete := models.SysMenu{Name: "删除菜单", Path: prefix + "/menu/delete", Method: datatypes.JSONSlice[string]{http.MethodDelete}, ParentID: menusMenu.ID, Type: models.MenuTypeButton}

	buttons := []*models.SysMenu{
		&userCreate, &userView, &userEdit, &userDelete,
		&roleCreate, &roleView, &roleEdit, &roleDelete, &roleAuth, &roleAuthList,
		&menuCreate, &menuView, &menuEdit, &menuDelete,
	}
	for _, b := range buttons {
		if err := db.Create(b).Error; err != nil {
			return err
		}
	}

	// ---------------- 创建角色 ----------------
	adminRole := models.SysRole{
		Name:  "admin",
		Code:  "admin",
		Menus: append(menus, buttons...),
	}
	editorRole := models.SysRole{
		Name:  "editor",
		Code:  "editor",
		Menus: []*models.SysMenu{&dashboard, &usersMenu, &rolesMenu, &userView, &userEdit, &roleView},
	}

	roles := []*models.SysRole{&adminRole, &editorRole}
	for _, r := range roles {
		if err := db.Create(r).Error; err != nil {
			return err
		}
	}

	// ---------------- 创建用户 ----------------
	adminUser := models.SysUser{Name: "admin", Email: "admin@example.com", Password: cryptox.HashMake("123456"), Roles: []*models.SysRole{&adminRole}}
	editorUser := models.SysUser{Name: "editor", Email: "editor@example.com", Password: cryptox.HashMake("123456"), Roles: []*models.SysRole{&editorRole}}

	users := []*models.SysUser{&adminUser, &editorUser}
	for _, u := range users {
		if err := db.Create(u).Error; err != nil {
			return err
		}
	}

	// ---------------- 同步 Casbin ----------------
	for _, r := range roles {
		if err := casbinx.SyncRole(e, r); err != nil {
			return err
		}
	}
	for _, u := range users {
		if err := casbinx.SyncUserRoles(e, u); err != nil {
			return err
		}
	}
	return nil
}
