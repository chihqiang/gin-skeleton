package config

// 导入必要的包
import (
	"wangzhiqiang/skeleton/pkg/database"
	"wangzhiqiang/skeleton/pkg/httpx"
	"wangzhiqiang/skeleton/pkg/logger"
	"wangzhiqiang/skeleton/pkg/redisx"
)

// Config 应用配置结构体
// 包含服务器、数据库和日志记录器的配置
type Config struct {
	Server   *httpx.Config    `yaml:"server"`   // 服务器配置
	Database *database.Config `yaml:"database"` // 数据库配置
	Logger   *logger.Config   `yaml:"logger"`   // 日志记录器配置
	Redis    *redisx.Config   `yaml:"redis"`
	Queue    *QueueConfig     `yaml:"queue"`
}

// QueueConfig 队列配置结构体
type QueueConfig struct {
	Redis *redisx.Config `yaml:"redis"`
	Queue string         `yaml:"queue"`
}
