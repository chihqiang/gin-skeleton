package mws

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const (
	defaultSessionName   = "go"
	defaultSessionSecret = "gin-skeleton"
)

type SessionConfig struct {
	Secret string `yaml:"secret" json:"secret"`
}

func Session(cfg *SessionConfig) gin.HandlerFunc {
	if cfg == nil {
		cfg = &SessionConfig{
			Secret: defaultSessionSecret,
		}
	}
	if cfg.Secret == "" {
		cfg.Secret = defaultSessionSecret
	}
	store := cookie.NewStore([]byte(cfg.Secret))
	return sessions.Sessions(defaultSessionName, store)
}
