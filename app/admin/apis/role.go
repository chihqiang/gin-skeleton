package apis

import (
	"context"
	"github.com/gin-gonic/gin"
	"wangzhiqiang/skeleton/app/admin/middlewares"
	"wangzhiqiang/skeleton/app/admin/service"
	"wangzhiqiang/skeleton/app/admin/types"
	"wangzhiqiang/skeleton/pkg/httpx"
)

type RoleApis struct {
	ctx     context.Context
	service *service.Service
}

func NewRole(ctx context.Context) *RoleApis {
	return &RoleApis{ctx: ctx, service: new(service.Service)}
}
func (r *RoleApis) List(c *gin.Context) {
	var req types.RoleListReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	resp, err := r.service.Role.List(r.ctx, &req)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, resp)
}

func (r *RoleApis) Create(c *gin.Context) {
	var req types.RoleReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	err := r.service.Role.Create(r.ctx, &req)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, req)
}

func (r *RoleApis) View(c *gin.Context) {
	var req types.RoleReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	info, err := r.service.Role.View(r.ctx, req.ID)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, info)
}

func (r *RoleApis) Edit(c *gin.Context) {
	var req types.RoleReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	err := r.service.Role.Edit(r.ctx, &req)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, map[string]string{})
}

func (r *RoleApis) Delete(c *gin.Context) {
	var req types.RoleReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	err := r.service.Role.Delete(r.ctx, req.ID)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, map[string]string{})
}

func (r *RoleApis) AuthList(c *gin.Context) {
	claims, err := middlewares.GetClaims(c)
	if err != nil {
		httpx.ApiNoAuth(c, err)
		return
	}
	menus, err := r.service.User.GetMenus(r.ctx, claims.UID)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, menus)
}

func (r *RoleApis) Auth(c *gin.Context) {
	var req types.RoleAuthReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	err := r.service.Role.Auth(r.ctx, &req)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, map[string]string{})
}
