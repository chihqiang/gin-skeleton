package types

import "wangzhiqiang/skeleton/pkg/database"

type UserListReq struct {
	database.PageRequest
	Email string `json:"email" form:"email" param:"email" uri:"email" query:"email"`
}

type UserReq struct {
	ID       uint   `json:"id,omitempty" form:"id" param:"id" uri:"id" query:"id"`
	Email    string `json:"email,omitempty" form:"email" param:"email" uri:"email" query:"email"`
	Name     string `json:"name,omitempty" form:"name" param:"name" uri:"name" query:"name"`
	Phone    string `json:"phone,omitempty" form:"phone" param:"phone" uri:"phone" query:"phone"`
	Password string `json:"password,omitempty" form:"password" param:"password" uri:"password" query:"password"`
	RoleIds  []uint `json:"role_ids,omitempty" form:"role_ids" param:"role_ids" uri:"role_ids" query:"role_ids"`
}
