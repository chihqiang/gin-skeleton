package app

import (
	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
	"wangzhiqiang/skeleton/config"
	"wangzhiqiang/skeleton/pkg/casbinx"
	"wangzhiqiang/skeleton/pkg/database"
	"wangzhiqiang/skeleton/pkg/httpx"
	"wangzhiqiang/skeleton/pkg/jwts"
	"wangzhiqiang/skeleton/pkg/logger"
	"wangzhiqiang/skeleton/pkg/queue"
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

func ProvideJWT(cfg *config.Config) *jwts.JWT {
	return jwts.NewJWT(cfg.JWT)
}

func ProvideQueue(db *gorm.DB, cfg *config.Config) (queue.IQueue, error) {
	if cfg.Queue == nil {
		cfg.Queue = &queue.Config{}
	}
	q, err := queue.New(cfg.Queue, db)
	if err != nil {
		return nil, err
	}
	tasks := queue.GetTasks()
	for _, task := range tasks {
		_ = q.Register(task)
	}
	return q, nil
}
