package mws

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Core() gin.HandlerFunc {
	return cors.New(cors.Config{
		// 允许的方法
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodHead,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Accept",            // 客户端可接收的内容类型
			"Accept-Encoding",   // 客户端支持的内容编码（gzip, deflate 等）
			"Accept-Language",   // 客户端首选语言
			"Authorization",     // 身份验证信息（如 Bearer token）
			"Content-Type",      // 请求体的内容类型
			"Content-Length",    // 请求体长度
			"Origin",            // 发起请求的源
			"Referer",           // 请求来源页面
			"User-Agent",        // 客户端信息
			"X-Requested-With",  // AJAX 请求标识
			"X-Forwarded-For",   // 代理转发的原始客户端 IP
			"X-Real-IP",         // 代理服务器提供的真实客户端 IP
			"X-Request-ID",      // 请求唯一 ID，用于追踪日志
			"X-Custom-Header",   // 示例自定义头
			"Cache-Control",     // 缓存控制
			"If-Modified-Since", // 条件 GET 请求头，协助缓存
			"If-None-Match",     // 条件 GET 请求头，协助缓存（ETag）
			"Range",             // 下载指定范围数据
			"Connection",        // 连接控制
			"Upgrade",           // 协议升级头
		},
		AllowAllOrigins: true,
		MaxAge:          12 * time.Hour,
	})
}
