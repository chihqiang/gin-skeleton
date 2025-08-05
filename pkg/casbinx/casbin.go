package casbinx

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

const (
	DefaultModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`
	DefaultTableName = "sys_casbin_rule"
)

type Config struct {
	DB        *gorm.DB
	Model     string `yaml:"model"`
	TableName string `yaml:"table_name"`
}

func New(cfg *Config) (*casbin.Enforcer, error) {
	if cfg.Model == "" {
		cfg.Model = DefaultModel
	}
	if cfg.TableName == "" {
		cfg.TableName = DefaultTableName
	}
	adapter, err := gormadapter.NewAdapterByDBUseTableName(cfg.DB, "", cfg.TableName)
	if err != nil {
		return nil, err
	}
	m, err := model.NewModelFromString(cfg.Model)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}
	// 自动加载策略
	if err := e.LoadPolicy(); err != nil {
		return nil, err
	}
	return e, nil
}

// IMenu 菜单接口
type IMenu interface {
	GetID() uint
	GetPath() string
	GetMethods() []string
}

// IRole 角色接口
type IRole interface {
	GetID() uint
	GetMenus() []IMenu
}

// IUser 用户接口
type IUser interface {
	GetID() uint
	GetRoles() []IRole
}

// SyncRole 同步单个角色的菜单权限到 Casbin（全量同步）
func SyncRole(e *casbin.Enforcer, role IRole) error {
	roleID := fmt.Sprintf("%d", role.GetID())
	// 获取已有策略
	oldPolicies, err := e.GetFilteredPolicy(0, roleID)
	if err != nil {
		return err
	}
	oldSet := make(map[string]struct{})
	for _, p := range oldPolicies {
		oldSet[p[1]+"#"+p[2]] = struct{}{}
	}
	// 构建新策略集合
	newSet := make(map[string]struct{})
	for _, menu := range role.GetMenus() {
		for _, m := range menu.GetMethods() {
			method := strings.ToUpper(strings.TrimSpace(m))
			if method == "" {
				method = http.MethodGet
			}
			key := menu.GetPath() + "#" + method
			newSet[key] = struct{}{}
			if _, exists := oldSet[key]; !exists {
				if _, err := e.AddPolicy(roleID, menu.GetPath(), method); err != nil {
					return err
				}
			}
		}
	}
	// 删除多余策略
	for _, p := range oldPolicies {
		key := p[1] + "#" + p[2]
		if _, exists := newSet[key]; !exists {
			if _, err := e.RemovePolicy(p); err != nil {
				return err
			}
		}
	}

	return nil
}

// SyncUserRoles 同步单个用户的角色关系到 Casbin（全量同步）
func SyncUserRoles(e *casbin.Enforcer, user IUser) error {
	userID := fmt.Sprintf("%d", user.GetID())
	// 获取已有分组绑定
	oldRoles, err := e.GetRolesForUser(userID)
	if err != nil {
		return err
	}
	oldSet := make(map[string]struct{})
	for _, r := range oldRoles {
		oldSet[r] = struct{}{}
	}
	// 构建新角色集合
	newSet := make(map[string]struct{})
	for _, role := range user.GetRoles() {
		roleID := fmt.Sprintf("%d", role.GetID())
		newSet[roleID] = struct{}{}
		if _, exists := oldSet[roleID]; !exists {
			if _, err := e.AddGroupingPolicy(userID, roleID); err != nil {
				return err
			}
		}
	}

	// 删除多余绑定
	for r := range oldSet {
		if _, exists := newSet[r]; !exists {
			if _, err := e.RemoveGroupingPolicy(userID, r); err != nil {
				return err
			}
		}
	}

	return nil
}

// CheckPermission 根据 http.Request 检查用户是否有权限
//
//	superUID 是超级管理员 可以无视权限
func CheckPermission(e *casbin.Enforcer, r *http.Request, userID uint) (bool, error) {
	// 获取请求路径和方法
	path := r.URL.Path
	method := strings.ToUpper(r.Method)
	if method == "" {
		method = http.MethodGet
	}
	return e.Enforce(fmt.Sprintf("%d", userID), path, method)
}

// GetUserMenuIDs 获取用户拥有权限的菜单 ID 列表
func GetUserMenuIDs(enforcer *casbin.Enforcer, uid uint) ([]uint, error) {
	userID := strconv.FormatUint(uint64(uid), 10)

	// 获取用户拥有的所有权限策略
	permissions, err := enforcer.GetFilteredPolicy(0, userID) // 第 0 列是 subject
	if err != nil {
		return nil, err
	}
	if len(permissions) == 0 {
		return nil, nil
	}

	menuIDMap := make(map[uint]struct{})
	for _, p := range permissions {
		if len(p) < 2 {
			continue
		}
		// p[1] 假设存的是菜单 ID
		mid, err := strconv.ParseUint(p[1], 10, 64)
		if err != nil {
			continue
		}
		menuIDMap[uint(mid)] = struct{}{}
	}

	// 转切片返回
	menuIDs := make([]uint, 0, len(menuIDMap))
	for id := range menuIDMap {
		menuIDs = append(menuIDs, id)
	}

	return menuIDs, nil
}

// SyncAll 批量同步用户、角色、菜单到 Casbin
func SyncAll(e *casbin.Enforcer, users []IUser) error {
	// 已同步过的角色缓存，避免重复同步
	syncedRoles := make(map[uint]struct{})

	for _, user := range users {
		// 先同步用户与角色绑定关系
		if err := SyncUserRoles(e, user); err != nil {
			return err
		}

		// 同步该用户的角色对应的菜单权限
		for _, role := range user.GetRoles() {
			roleID := role.GetID()
			if _, exists := syncedRoles[roleID]; exists {
				continue // 角色已同步过，跳过
			}
			if err := SyncRole(e, role); err != nil {
				return err
			}
			syncedRoles[roleID] = struct{}{}
		}
	}
	return nil
}
