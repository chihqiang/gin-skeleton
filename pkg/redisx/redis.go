package redisx

import "github.com/redis/go-redis/v9"

type Config struct {
	Addr     string `yaml:"addr"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func New(cfg *Config) *redis.Client {
	opt := &redis.Options{}
	if cfg.Addr != "" {
		opt.Addr = cfg.Addr
	}
	if cfg.Username != "" {
		opt.Username = cfg.Username
	}
	if cfg.Password != "" {
		opt.Password = cfg.Password
	}
	if cfg.DB >= 0 {
		opt.DB = cfg.DB
	}
	return redis.NewClient(opt)
}
