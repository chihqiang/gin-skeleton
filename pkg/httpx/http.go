package httpx

// 导入必要的包
import (
	"wangzhiqiang/skeleton/pkg/httpx/middlewares" // 自定义中间件

	"github.com/gin-contrib/gzip" // Gzip 压缩中间件
	"github.com/gin-gonic/gin"    // Gin Web 框架
)

// Config HTTP 服务器配置结构体
// 包含服务器运行所需的各种配置参数
type Config struct {
	Port         int                        `yaml:"port"`          // 监听端口号
	Mode         string                     `yaml:"mode"`          // 运行模式（如：debug、release）
	ReadTimeout  int                        `yaml:"read_timeout"`  // 读取请求的超时时间（单位：秒）
	WriteTimeout int                        `yaml:"write_timeout"` // 写响应的超时时间（单位：秒）
	IdleTimeout  int                        `yaml:"idle_timeout"`  // 空闲连接的超时时间（单位：秒）
	Session      *middlewares.SessionConfig `yaml:"session"`       // 会话配置
}

// IController 控制器接口
// 所有控制器都需要实现这个接口来注册路由
type IController interface {
	// Routes 注册控制器的路由
	// 参数 g: Gin 路由接口
	Routes(g gin.IRoutes)
}

// Router 路由器结构体
// 包含多个控制器，用于统一注册路由
type Router struct {
	Controllers []IController // 控制器列表
}

// Route 创建并配置 Gin 引擎
// 参数 cfg: HTTP 服务器配置
// 返回值: 配置好的 Gin 引擎实例
func (route *Router) HTTP(cfg *Config) *gin.Engine {
	// 设置 Gin 运行模式
	gin.SetMode(cfg.Mode)

	// 创建新的 Gin 引擎
	engine := gin.New()

	// 使用恢复中间件，捕获并处理 panic
	engine.Use(gin.Recovery())

	// 使用日志中间件
	engine.Use(gin.Logger())

	// 使用安全中间件
	engine.Use(middlewares.Security())

	// 使用会话中间件
	engine.Use(middlewares.Session(cfg.Session))

	// 使用 Gzip 压缩中间件
	engine.Use(gzip.Gzip(gzip.DefaultCompression))

	// 注册所有控制器的路由
	for _, controller := range route.Controllers {
		if controller == nil {
			continue
		}
		controller.Routes(engine)
	}

	return engine
}
