package types

import "wangzhiqiang/skeleton/app/admin/models"

type LoginReq struct {
	Email    string `json:"email" form:"email" param:"email" uri:"email" query:"email"`
	Password string `json:"password" form:"password" param:"password" uri:"password" query:"password"`
	IP       string
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token" param:"refresh_token" uri:"refresh_token" query:"refresh_token" binding:"required"`
}

type LoginResp struct {
	Token        string            `json:"token"`
	RefreshToken string            `json:"refresh_token"`
	ExpiresIn    int               `json:"expires_in"`
	Menus        []*models.SysMenu `json:"menus"`
}
