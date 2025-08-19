package apis

import (
	"context"
	"github.com/gin-gonic/gin"
	"wangzhiqiang/skeleton/app/admin/service"
	"wangzhiqiang/skeleton/app/admin/types"
	"wangzhiqiang/skeleton/pkg/httpx"
)

type AuthApis struct {
	ctx     context.Context
	service *service.Service
}

func NewAuth(ctx context.Context) *AuthApis {
	return &AuthApis{ctx: ctx, service: new(service.Service)}
}

func (l *AuthApis) Login(c *gin.Context) {
	var req types.LoginReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	req.IP = c.ClientIP()
	resp, err := l.service.Auth.Login(l.ctx, &req)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess[*types.LoginResp](c, resp)
}

func (l *AuthApis) Refresh(c *gin.Context) {
	var req types.RefreshTokenReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	resp, err := l.service.Auth.Refresh(l.ctx, &req)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess[*types.LoginResp](c, resp)
}
