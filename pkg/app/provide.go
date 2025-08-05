package app

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"wangzhiqiang/skeleton/config"
	"wangzhiqiang/skeleton/pkg/database"
	"wangzhiqiang/skeleton/pkg/httpx"
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

func ProvideHTTPServer(router *httpx.Router, cfg *config.Config) http.Handler {
	return router.HTTP(cfg.Server)
}

func ProvideRouter(controllers []httpx.IController) *httpx.Router {
	return &httpx.Router{Controllers: controllers}
}

func ProvideControllers() []httpx.IController {
	return httpx.GetControllers()
}

func ProvideRedis(cfg *config.Config) *redis.Client {
	return redisx.New(cfg.Redis)
}

func ProvideQueueTask() []queue.ITask {
	return queue.GetTasks()
}
func ProvideQueue(redis *redis.Client, cfg *config.Config, tasks []queue.ITask) queue.IQueue {
	if cfg.Queue == nil {
		cfg.Queue = &config.QueueConfig{Queue: "skeleton"}
	}
	if cfg.Queue.Redis != nil {
		redis = redisx.New(cfg.Queue.Redis)
	}
	redisQueue := queue.NewRedisQueue(redis, cfg.Queue.Queue)
	for _, task := range tasks {
		_ = redisQueue.Register(task)
	}
	return redisQueue
}
