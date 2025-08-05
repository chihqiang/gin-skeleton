package types

import "wangzhiqiang/skeleton/pkg/database"

type RoleListReq struct {
	database.PageRequest
}

type RoleReq struct {
	ID      uint   `json:"id,omitempty" form:"id" param:"id" uri:"id" query:"id"` // 更新时用
	Name    string `json:"name,omitempty" form:"name" param:"name" uri:"name" query:"name"`
	Code    string `json:"code,omitempty" form:"code" param:"code" uri:"code" query:"code"`
	Remark  string `json:"remark,omitempty" form:"remark" param:"remark" uri:"remark" query:"remark"`
	MenuIds []uint `json:"menu_ids,omitempty" form:"menu_ids" param:"menu_ids" uri:"menu_ids" query:"menu_ids"`
}

type RoleAuthReq struct {
	RoleID  uint   `json:"role_id" form:"role_id"`
	MenuIDs []uint `json:"menu_ids" form:"menu_ids"`
}
