package models

import "time"

// SysAccessLog 访问日志
type SysAccessLog struct {
	ID uint `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`

	// ---------------------- 请求 ----------------------
	UserID    uint   `gorm:"index;type:bigint;comment:用户ID" json:"user_id"`
	Path      string `gorm:"index;type:varchar(255);comment:请求路径" json:"path"`
	Method    string `gorm:"type:varchar(10);comment:请求方法" json:"method"`
	Request   string `gorm:"type:longtext;comment:请求内容(截断1000字符)" json:"request"`
	Ip        string `gorm:"type:varchar(45);comment:请求IP" json:"ip"`
	RequestID string `gorm:"type:varchar(100);comment:请求唯一表示" json:"request_id"`
	UserAgent string `gorm:"type:varchar(255);comment:请求User-Agent" json:"user_agent"`
	// ---------------------- 响应 ----------------------
	Status   int    `gorm:"type:int;index;comment:响应状态" json:"status"`
	Latency  int64  `gorm:"type:bigint;comment:延迟(毫秒)" json:"latency"` // 存储毫秒
	Response string `gorm:"type:longtext;comment:响应内容(截断1000字符)" json:"response"`

	CreatedAt time.Time `gorm:"autoCreateTime;index;comment:创建时间" json:"created_at"`
}
