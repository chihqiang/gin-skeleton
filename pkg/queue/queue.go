package queue

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
	"wangzhiqiang/skeleton/pkg/database"
	"wangzhiqiang/skeleton/pkg/redisx"
)

type IQueue interface {
	// Register 注册任务，比如初始化或者将任务写入某个存储（DB、Redis、内存等）
	Register(task ITask) error

	// Push 推送任务，可以带延时（适合定时任务、延迟队列）
	Push(task ITask, delay time.Duration) error

	// Start 启动队列监听，interval 表示检查/拉取任务的间隔
	Start(ctx context.Context, interval time.Duration)

	// Stop 停止队列
	Stop()
}

type Config struct {
	Type  string           `yaml:"type" json:"type,omitempty"`
	Queue string           `yaml:"queue" json:"queue,omitempty"`
	DB    *database.Config `yaml:"db" json:"db,omitempty"`
	Redis *redisx.Config   `yaml:"redis" json:"redis,omitempty"`
}

func New(cfg *Config, db *gorm.DB, rds *redis.Client) (IQueue, error) {
	if cfg.Type == "" {
		cfg.Type = "db"
	}
	if cfg.Queue == "" {
		cfg.Queue = "skeleton"
	}
	var err error
	switch cfg.Type {

	case "redis":
		if cfg.Redis != nil {
			rds = redisx.New(cfg.Redis)
		}
		return NewRedisQueue(rds, cfg.Queue), err
	default:
		if cfg.DB != nil {
			db, err = database.Init(cfg.DB)
		}
		return NewGormQueue(db), err
	}
}
