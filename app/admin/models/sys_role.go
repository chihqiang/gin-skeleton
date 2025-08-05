package models

import (
	"wangzhiqiang/skeleton/pkg/casbinx"
	"wangzhiqiang/skeleton/pkg/database"
)

type SysRole struct {
	database.BaseModel
	Name   string     `gorm:"type:varchar(50);uniqueIndex;comment:角色名称" json:"name"`
	Code   string     `gorm:"type:varchar(50);uniqueIndex;comment:角色编码" json:"code"`
	Remark string     `gorm:"type:varchar(255);comment:备注" json:"remark"`
	Users  []*SysUser `gorm:"many2many:sys_user_roles;" json:"users"`
	Menus  []*SysMenu `gorm:"many2many:sys_role_menus;" json:"menus"`
}

// GetID 返回角色ID（string 类型）
func (r *SysRole) GetID() uint {
	return r.ID
}

// GetMenus 返回角色绑定的菜单
func (r *SysRole) GetMenus() []casbinx.IMenu {
	menus := make([]casbinx.IMenu, len(r.Menus))
	for i, m := range r.Menus {
		menus[i] = m // SysMenu 必须实现 IMenu
	}
	return menus
}
