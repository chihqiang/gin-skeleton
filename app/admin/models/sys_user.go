package models

import (
	"time"
	"wangzhiqiang/skeleton/pkg/casbinx"
	"wangzhiqiang/skeleton/pkg/database"
)

type SysUser struct {
	database.BaseModel
	Email     string    `gorm:"type:varchar(100);uniqueIndex;default:'';index:idx_email;comment:电子邮箱" json:"email"`
	Name      string    `gorm:"type:varchar(50);default:'';comment:昵称" json:"name"`
	Phone     string    `gorm:"type:bigint unsigned;default:0;comment:手机号" json:"phone"`
	Password  string    `gorm:"type:varchar(255);default:'';comment:密码" json:"-"`
	LastLogin time.Time `gorm:"default:null;comment:最后登录时间" json:"last_login,omitempty"`
	LastIp    string    `gorm:"type:varchar(100);default:'';comment:最后登录IP" json:"last_ip,omitempty"`

	Roles []*SysRole `gorm:"many2many:sys_user_roles;" json:"roles"`
}

// GetID 返回用户ID（string 类型）
func (u *SysUser) GetID() uint {
	return u.ID
}
func (u *SysUser) GetName() string {
	return u.Name
}

// GetRoles 返回用户关联的角色列表
func (u *SysUser) GetRoles() []casbinx.IRole {
	roles := make([]casbinx.IRole, len(u.Roles))
	for i, r := range u.Roles {
		roles[i] = r // SysRole 必须实现 IRole
	}
	return roles
}
