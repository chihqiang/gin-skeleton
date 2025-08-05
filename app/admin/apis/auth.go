package apis

import (
	"context"
	"github.com/gin-gonic/gin"
	"wangzhiqiang/skeleton/app/admin/service"
	"wangzhiqiang/skeleton/app/admin/types"
	"wangzhiqiang/skeleton/pkg/httpx"
)

type Auth struct {
	ctx context.Context
}

func NewAuth(ctx context.Context) *Auth {
	return &Auth{ctx: ctx}
}

func (l *Auth) Login(c *gin.Context) {
	var req types.LoginReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	req.IP = c.ClientIP()
	resp, err := new(service.AuthService).Login(l.ctx, &req)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess[*types.LoginResp](c, resp)
}

func (l *Auth) Refresh(c *gin.Context) {
	var req types.RefreshTokenReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	resp, err := new(service.AuthService).Refresh(l.ctx, &req)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess[*types.LoginResp](c, resp)
}
