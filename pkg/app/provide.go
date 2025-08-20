package app

import (
	"github.com/casbin/casbin/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"wangzhiqiang/skeleton/config"
	"wangzhiqiang/skeleton/pkg/casbinx"
	"wangzhiqiang/skeleton/pkg/database"
	"wangzhiqiang/skeleton/pkg/httpx"
	"wangzhiqiang/skeleton/pkg/jwts"
	"wangzhiqiang/skeleton/pkg/logger"
	"wangzhiqiang/skeleton/pkg/queue"
	"wangzhiqiang/skeleton/pkg/redisx"
)

func ProvideLogger(cfg *config.Config) (logger.ILogger, error) {
	return logger.Init(cfg.Logger)
}

func ProvideDatabase(cfg *config.Config) (*gorm.DB, error) {
	return database.Init(cfg.Database)
}

func ProvideEnforcer(db *gorm.DB) (*casbin.Enforcer, error) {
	return casbinx.New(&casbinx.Config{DB: db})
}

func ProvideHTTPServer(cfg *config.Config) *httpx.HTTP {
	return httpx.NewHTTP(cfg.Server)
}

func ProvideRedis(cfg *config.Config) *redis.Client {
	return redisx.New(cfg.Redis)
}

func ProvideQueueTask() []queue.ITask {
	return queue.GetTasks()
}

func ProvideJWT(cfg *config.Config) *jwts.JWT {
	return jwts.NewJWT(cfg.JWT)
}

func ProvideQueue(rds *redis.Client, db *gorm.DB, cfg *config.Config, tasks []queue.ITask) (queue.IQueue, error) {
	if cfg.Queue == nil {
		cfg.Queue = &queue.Config{Queue: "skeleton"}
	}
	q, err := queue.New(cfg.Queue, db, rds)
	if err != nil {
		return nil, err
	}
	for _, task := range tasks {
		_ = q.Register(task)
	}
	return q, nil
}
