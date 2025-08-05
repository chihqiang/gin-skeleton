package apis

import (
	"context"
	"github.com/gin-gonic/gin"
	"wangzhiqiang/skeleton/app/admin/middlewares"
	"wangzhiqiang/skeleton/app/admin/service"
	"wangzhiqiang/skeleton/app/admin/types"
	"wangzhiqiang/skeleton/pkg/httpx"
)

type Role struct {
	ctx context.Context
}

func NewRole(ctx context.Context) *Role {
	return &Role{ctx: ctx}
}
func (r *Role) List(c *gin.Context) {
	var req types.RoleListReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	resp, err := new(service.RoleService).List(r.ctx, &req)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, resp)
}

func (r *Role) Create(c *gin.Context) {
	var req types.RoleReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	err := new(service.RoleService).Create(r.ctx, &req)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, req)
}

func (r *Role) View(c *gin.Context) {
	var req types.RoleReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	info, err := new(service.RoleService).View(r.ctx, req.ID)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, info)
}

func (r *Role) Edit(c *gin.Context) {
	var req types.RoleReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	err := new(service.RoleService).Edit(r.ctx, &req)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, map[string]string{})
}

func (r *Role) Delete(c *gin.Context) {
	var req types.RoleReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	err := new(service.RoleService).Delete(r.ctx, req.ID)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, map[string]string{})
}

func (r *Role) AuthList(c *gin.Context) {
	claims, err := middlewares.GetClaims(c)
	if err != nil {
		httpx.ApiNoAuth(c, err)
		return
	}
	menus, err := new(service.UserService).GetMenus(r.ctx, claims.UID)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, menus)
}

func (r *Role) Auth(c *gin.Context) {
	var req types.RoleAuthReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	err := new(service.RoleService).Auth(r.ctx, &req)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, map[string]string{})
}
