package httpx

import (
	"context"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
	"wangzhiqiang/skeleton/pkg/httpx/mws"
)

// Config HTTP 服务器配置结构体
// 包含服务器运行所需的各种配置参数
type Config struct {
	Port         int                `yaml:"port" json:"port,omitempty"`                   // 监听端口号
	Mode         string             `yaml:"mode" json:"mode,omitempty"`                   // 运行模式（如：debug、release）
	ReadTimeout  int                `yaml:"read_timeout" json:"read_timeout,omitempty"`   // 读取请求的超时时间（单位：秒）
	WriteTimeout int                `yaml:"write_timeout" json:"write_timeout,omitempty"` // 写响应的超时时间（单位：秒）
	IdleTimeout  int                `yaml:"idle_timeout" json:"idle_timeout,omitempty"`   // 空闲连接的超时时间（单位：秒）
	Session      *mws.SessionConfig `yaml:"session" json:"session,omitempty"`             // 会话配置
}

type HTTP struct {
	config     *Config
	serverLock sync.Mutex
	httpServer *http.Server
}

func NewHTTP(cfg *Config) *HTTP {
	return &HTTP{config: cfg}
}

func (h *HTTP) newGin(ctx context.Context) *gin.Engine {
	// 设置 Gin 运行模式
	gin.SetMode(h.config.Mode)
	// 创建新的 Gin 引擎
	engine := gin.Default()
	// 使用安全中间件
	engine.Use(mws.Security())
	//request
	engine.Use(mws.RequestID())
	// 使用会话中间件
	engine.Use(mws.Session(h.config.Session))
	// 使用 Gzip 压缩中间件
	engine.Use(gzip.Gzip(gzip.DefaultCompression))
	// 注册所有控制器的路由
	for _, route := range routes {
		if route == nil {
			continue
		}
		if err := route.Routes(ctx, engine); err != nil {
			panic(fmt.Sprintf("register route failed: %v", err))
		}
	}
	return engine
}

func (h *HTTP) Start(ctx context.Context) error {
	handler := h.newGin(ctx)
	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", h.config.Port),           // 服务器地址
		Handler:      handler,                                            // 处理器
		ReadTimeout:  time.Duration(h.config.ReadTimeout) * time.Second,  // 读取超时
		WriteTimeout: time.Duration(h.config.WriteTimeout) * time.Second, // 写入超时
		IdleTimeout:  time.Duration(h.config.IdleTimeout) * time.Second,  // 空闲超时
	}
	h.serverLock.Lock()
	h.httpServer = srv
	h.serverLock.Unlock()
	return srv.ListenAndServe()
}

func (h *HTTP) Stop(ctx context.Context) error {
	return h.httpServer.Shutdown(ctx)
}
