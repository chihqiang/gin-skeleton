package models

import "time"

// SysAccessLog AccessLog 用于保存请求日志到数据库的结构体
type SysAccessLog struct {
	ID            uint      `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	UserID        uint      `gorm:"index;type:bigint;comment:用户ID" json:"user_id"`
	Path          string    `gorm:"type:varchar(255);comment:请求路径" json:"path"`
	Method        string    `gorm:"type:varchar(10);comment:请求方法" json:"method"`
	Body          string    `gorm:"type:longtext;comment:请求内容" json:"body"` // 超大请求内容
	IP            string    `gorm:"type:varchar(45);comment:请求IP" json:"ip"`
	UserAgent     string    `gorm:"type:varchar(255);comment:请求User-Agent" json:"user_agent"`
	ProcessTimeMs int64     `gorm:"comment:请求处理耗时毫秒" json:"process_time_ms"`
	CreatedAt     time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
}
