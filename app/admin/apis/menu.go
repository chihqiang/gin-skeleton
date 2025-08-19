package apis

import (
	"context"
	"github.com/gin-gonic/gin"
	"wangzhiqiang/skeleton/app/admin/service"
	"wangzhiqiang/skeleton/app/admin/types"
	"wangzhiqiang/skeleton/pkg/httpx"
)

type MenuApis struct {
	ctx     context.Context
	service *service.Service
}

func NewMenu(ctx context.Context) *MenuApis {
	return &MenuApis{ctx: ctx, service: new(service.Service)}
}
func (m *MenuApis) List(c *gin.Context) {
	var req types.MenuListReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	resp, err := m.service.Menu.List(m.ctx, &req)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, resp)
}

func (m *MenuApis) Create(c *gin.Context) {
	var req types.MenuReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	if err := m.service.Menu.Create(m.ctx, &req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, map[string]string{})
}

func (m *MenuApis) View(c *gin.Context) {
	var req types.IDReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	info, err := m.service.Menu.Get(m.ctx, req.ID)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, info)
}

func (m *MenuApis) Edit(c *gin.Context) {
	var req types.MenuReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	err := m.service.Menu.Edit(m.ctx, &req)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, map[string]string{})
}

func (m *MenuApis) Delete(c *gin.Context) {
	var req types.IDReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	err := m.service.Menu.Delete(m.ctx, req.ID)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, map[string]string{})
}
