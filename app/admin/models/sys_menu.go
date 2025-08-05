package models

import (
	"gorm.io/datatypes"
	"net/http"
	"wangzhiqiang/skeleton/pkg/database"
)

type MenuType string

const (
	MenuTypeMenu   MenuType = "menu"   // 普通菜单
	MenuTypeButton MenuType = "button" // 按钮/操作权限
)

type SysMenu struct {
	database.BaseModel
	ParentID  uint                        `gorm:"default:0;comment:父级菜单ID" json:"parent_id"`
	Name      string                      `gorm:"type:varchar(50);comment:菜单名称" json:"name"`
	Method    datatypes.JSONSlice[string] `gorm:"type:varchar(50);default:'GET';comment:请求方法'" json:"method"`
	Path      string                      `gorm:"type:varchar(100);comment:路由路径" json:"path"`
	Component string                      `gorm:"type:varchar(100);comment:前端组件路径" json:"component"`
	Icon      string                      `gorm:"type:varchar(50);comment:菜单图标" json:"icon"`
	Sort      int                         `gorm:"default:0;comment:排序" json:"sort"`
	Hidden    bool                        `gorm:"default:false;comment:是否隐藏" json:"hidden"`
	Type      MenuType                    `gorm:"type:varchar(10);default:'menu';comment:类型（menu=菜单, button=按钮）" json:"type"`
	Roles     []SysRole                   `gorm:"many2many:sys_role_menus;" json:"roles"`
	Children  []*SysMenu                  `gorm:"-" json:"children,omitempty"`

	Selected bool `gorm:"-" json:"selected,omitempty"`
}

func (m *SysMenu) GetID() uint {
	return m.ID
}

func (m *SysMenu) GetParentID() uint {
	return m.ParentID
}

func (m *SysMenu) SetChildren(children []*SysMenu) {
	m.Children = children
}

func (m *SysMenu) GetSort() int {
	return m.Sort
}

// GetPath 返回菜单路径
func (m *SysMenu) GetPath() string {
	return m.Path
}

// GetMethods 返回菜单绑定的请求方法（切片）
func (m *SysMenu) GetMethods() []string {
	if len(m.Method) == 0 {
		// 默认 GET 方法
		return []string{http.MethodGet}
	}
	return m.Method
}

func (m *SysMenu) SetSelected(selected bool) {
	m.Selected = selected
}
