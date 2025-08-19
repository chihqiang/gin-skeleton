package apis

import (
	"context"
	"github.com/gin-gonic/gin"
	"wangzhiqiang/skeleton/app/admin/service"
	"wangzhiqiang/skeleton/app/admin/types"
	"wangzhiqiang/skeleton/pkg/httpx"
)

type UserApis struct {
	ctx     context.Context
	service *service.Service
}

func NewUser(ctx context.Context) *UserApis {
	return &UserApis{ctx: ctx, service: new(service.Service)}
}

func (u *UserApis) List(c *gin.Context) {
	var req types.UserListReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	resp, err := u.service.User.List(u.ctx, &req)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, resp)
}

func (u *UserApis) Create(c *gin.Context) {
	var req types.UserReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	if err := u.service.User.Create(u.ctx, &req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, map[string]string{})
}

func (u *UserApis) View(c *gin.Context) {
	var req types.IDReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	info, err := u.service.User.Get(u.ctx, req.ID)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, info)
}

func (u *UserApis) Edit(c *gin.Context) {
	var req types.UserReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	if err := u.service.User.Edit(u.ctx, &req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, map[string]string{})
}

func (u *UserApis) Delete(c *gin.Context) {
	var req types.IDReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	err := u.service.User.Delete(u.ctx, &req)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, map[string]string{})
}
