package apis

import (
	"context"
	"github.com/gin-gonic/gin"
	"wangzhiqiang/skeleton/app/admin/service"
	"wangzhiqiang/skeleton/app/admin/types"
	"wangzhiqiang/skeleton/pkg/httpx"
)

type Menu struct {
	ctx context.Context
}

func NewMenu(ctx context.Context) *Menu {
	return &Menu{ctx: ctx}
}
func (m *Menu) List(c *gin.Context) {
	var req types.MenuListReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	resp, err := new(service.MenuService).List(m.ctx, &req)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, resp)
}

func (m *Menu) Create(c *gin.Context) {
	var req types.MenuReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	if err := new(service.MenuService).Create(m.ctx, &req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, map[string]string{})
}

func (m *Menu) View(c *gin.Context) {
	var req types.IDReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	info, err := new(service.MenuService).Get(m.ctx, req.ID)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, info)
}

func (m *Menu) Edit(c *gin.Context) {
	var req types.MenuReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	err := new(service.MenuService).Edit(m.ctx, &req)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, map[string]string{})
}

func (m *Menu) Delete(c *gin.Context) {
	var req types.IDReq
	if err := c.ShouldBind(&req); err != nil {
		httpx.ApiError(c, err)
		return
	}
	err := new(service.MenuService).Delete(m.ctx, req.ID)
	if err != nil {
		httpx.ApiError(c, err)
		return
	}
	httpx.ApiSuccess(c, map[string]string{})
}
