package config

// 导入必要的包
import (
	"gopkg.in/yaml.v3"
	"os"
	"wangzhiqiang/skeleton/pkg/database"
	"wangzhiqiang/skeleton/pkg/httpx"
	"wangzhiqiang/skeleton/pkg/httpx/mws"
	"wangzhiqiang/skeleton/pkg/jwts"
	"wangzhiqiang/skeleton/pkg/logger"
	"wangzhiqiang/skeleton/pkg/queue"
)

// Config 应用配置结构体
// 包含服务器、数据库和日志记录器的配置
type Config struct {
	Server   *httpx.Config    `yaml:"server" json:"server,omitempty"`     // 服务器配置，包括监听端口、运行模式、超时设置及会话配置
	Database *database.Config `yaml:"database" json:"database,omitempty"` // 数据库配置，包括驱动类型、连接信息、连接池参数等
	Logger   *logger.Config   `yaml:"logger" json:"logger,omitempty"`     // 日志记录器配置，包括日志级别、文件路径、格式、切割与压缩策略
	Queue    *queue.Config    `yaml:"queue" json:"queue,omitempty"`       // 队列配置
	JWT      *jwts.Config     `yaml:"jwt" json:"jwt,omitempty"`           // JWT 配置，包括密钥、过期时间、签发者和受众信息
	System   *SystemConfig    `yaml:"system" json:"system,omitempty"`     // 系统配置，包括超级管理员 ID 等全局系统参数
}

type SystemConfig struct {
	SuperAdminUID uint `yaml:"super_admin_uid" json:"super_admin_uid,omitempty"`
}

var (
	defaultDB = &database.Config{
		Driver: "sqlite",
		DBName: "runtime/skeleton.db",
	}
	defaultServerSession = &mws.SessionConfig{
		Secret: "gin-skeleton",
	}
	defaultLogger = &logger.Config{
		Level: "debug",
		//Path:       "runtime/log/app.log",
		Path:       "",
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
		Format:     "json",
	}
	defaultQueue  = &queue.Config{}
	defaultSystem = &SystemConfig{
		SuperAdminUID: 1,
	}
	defaultJWT = &jwts.Config{
		Secret: "vosMykI4axI9IrUuI8JYxlaHnnEWLvfNrWE3gOwOBBk=",
	}
)

// Load 加载配置
func Load(path string) (*Config, error) {
	cfg := &Config{}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	// 设置默认值，如果子配置为 nil 则初始化
	if cfg.System == nil {
		cfg.System = defaultSystem
	}
	if cfg.Server.Session == nil {
		cfg.Server.Session = defaultServerSession
	}
	if cfg.Database == nil {
		cfg.Database = defaultDB
	}
	// 初始化指针字段，保证使用时不为 nil
	if cfg.Logger == nil {
		cfg.Logger = defaultLogger
	}
	if cfg.Queue == nil {
		cfg.Queue = defaultQueue
	}
	if cfg.JWT == nil {
		cfg.JWT = defaultJWT
	}
	return cfg, nil
}
