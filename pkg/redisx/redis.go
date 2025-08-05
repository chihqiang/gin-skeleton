package redisx

import "github.com/redis/go-redis/v9"

type Config struct {
	Addr     string `yaml:"addr" json:"addr,omitempty"`
	Username string `yaml:"username" json:"username,omitempty"`
	Password string `yaml:"password" json:"password,omitempty"`
	DB       int    `yaml:"db" json:"db,omitempty"`
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
