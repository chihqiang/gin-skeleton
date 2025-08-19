package service

import (
	"context"
	"fmt"
	"time"
	"wangzhiqiang/skeleton/app/admin/models"
	"wangzhiqiang/skeleton/app/admin/types"
	"wangzhiqiang/skeleton/pkg/app"
	"wangzhiqiang/skeleton/pkg/cryptox"
)

type AuthService struct {
}

func (s *AuthService) Login(ctx context.Context, req *types.LoginReq) (*types.LoginResp, error) {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return nil, err
	}
	//查询用户信息
	var user models.SysUser
	db := apps.DB
	if err := db.Model(models.SysUser{}).Where(models.SysUser{Email: req.Email}).First(&user).Error; err != nil {
		return nil, fmt.Errorf("邮箱不存在")
	}
	// 验证密码
	if !cryptox.HashVerify(req.Password, user.Password) {
		return nil, fmt.Errorf("密码错误")
	}
	jwt := apps.JWT
	//生成JWT token
	token, err := jwt.BuildAccessToken(&user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := jwt.BuildRefreshToken(&user)
	if err != nil {
		return nil, err
	}
	// 更新用户登录信息
	db.Model(models.SysUser{}).Where("id = ?", user.ID).Updates(&models.SysUser{
		LastLogin: time.Now(),
		LastIp:    req.IP,
	})
	//获取当前用户菜单权限
	menus, _ := new(UserService).GetUserMenus(ctx, user.ID)
	return &types.LoginResp{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresIn:    jwt.GetAccessExpiresIn(),
		Menus:        menus,
	}, nil
}
func (s *AuthService) Refresh(ctx context.Context, req *types.RefreshTokenReq) (*types.LoginResp, error) {
	apps, err := app.GetApps(ctx)
	if err != nil {
		return nil, err
	}
	jwts := apps.JWT
	newAccess, newRefresh, err := jwts.Refresh(req.RefreshToken)
	if err != nil {
		return nil, err
	}
	// 返回新的 token
	return &types.LoginResp{
		Token:        newAccess,
		RefreshToken: newRefresh,
		ExpiresIn:    jwts.GetAccessExpiresIn(),
	}, nil
}
