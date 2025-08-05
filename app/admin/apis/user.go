package apis

import (
	"context"
	"github.com/gin-gonic/gin"
	"wangzhiqiang/skeleton/app/admin/service"
	"wangzhiqiang/skeleton/app/admin/types"
	"wangzhiqiang/skeleton/pkg/httpx"
)

type User struct {
	ctx context.Context
}

func NewUser(ctx context.Context) *User {
	return &User{ctx: ctx}
}

func (u *User) List(c *gin.Context) {
	var req types.UserListReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	resp, err := new(service.UserService).List(u.ctx, &req)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, resp)
}

func (u *User) Create(c *gin.Context) {
	var req types.UserReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	if err := new(service.UserService).Create(u.ctx, &req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, map[string]string{})
}

func (u *User) View(c *gin.Context) {
	var req types.IDReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	info, err := new(service.UserService).Get(u.ctx, req.ID)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, info)
}

func (u *User) Edit(c *gin.Context) {
	var req types.UserReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	if err := new(service.UserService).Edit(u.ctx, &req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, map[string]string{})
}

func (u *User) Delete(c *gin.Context) {
	var req types.IDReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	err := new(service.UserService).Delete(u.ctx, &req)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, map[string]string{})
}
