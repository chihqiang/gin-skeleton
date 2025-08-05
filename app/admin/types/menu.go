package types

import "wangzhiqiang/skeleton/pkg/database"

type MenuListReq struct {
	database.PageRequest
}

// MenuReq 用于创建或更新菜单
type MenuReq struct {
	ID        uint     `json:"id,omitempty" form:"id" param:"id" uri:"id" query:"id"` // 更新时用
	ParentID  uint     `json:"parent_id,omitempty" form:"parent_id" param:"parent_id" uri:"parent_id" query:"parent_id"`
	Name      string   `json:"name,omitempty" form:"name" param:"name" uri:"name" query:"name" binding:"required"`
	Method    []string `json:"method,omitempty" form:"method" param:"method" uri:"method" query:"method"`
	Path      string   `json:"path,omitempty" form:"path" param:"path" uri:"path" query:"path" binding:"required"`
	Component string   `json:"component,omitempty" form:"component" param:"component" uri:"component" query:"component"`
	Icon      string   `json:"icon,omitempty" form:"icon" param:"icon" uri:"icon" query:"icon"`
	Sort      int      `json:"sort,omitempty" form:"sort" param:"sort" uri:"sort" query:"sort"`
	Hidden    bool     `json:"hidden,omitempty" form:"hidden" param:"hidden" uri:"hidden" query:"hidden"`
	Type      string   `json:"type,omitempty" form:"type" param:"type" uri:"type" query:"type" binding:"oneof=menu button"`
	RoleIds   []uint   `json:"role_ids,omitempty" form:"role_ids" param:"role_ids" uri:"role_ids" query:"role_ids"`
}
