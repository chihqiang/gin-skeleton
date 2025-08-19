package middlewares

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"wangzhiqiang/skeleton/app/admin/middlewares"
	"wangzhiqiang/skeleton/app/models"
	"wangzhiqiang/skeleton/pkg/httpx/mws"
)

func AccessLog(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			body   []byte
			userID uint
			err    error
		)
		// 捕获请求 body
		switch c.Request.Method {
		case http.MethodGet:
			query, _ := url.QueryUnescape(c.Request.URL.RawQuery)
			split := strings.Split(query, "&")
			m := make(map[string]string)
			for _, v := range split {
				kv := strings.SplitN(v, "=", 2)
				if len(kv) == 2 {
					m[kv[0]] = kv[1]
				}
			}
			body, _ = json.Marshal(&m)
		default:
			body, err = io.ReadAll(c.Request.Body)
			if err == nil {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		}
		// 获取用户 ID
		claims, _ := middlewares.GetClaims(c)
		if claims != nil && claims.UID != 0 {
			userID = claims.UID
		}
		// 捕获响应 body
		respBody := &bytes.Buffer{}
		writer := &bodyWriter{ResponseWriter: c.Writer, body: respBody}
		c.Writer = writer
		// 记录开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 构建日志记录
		record := models.SysAccessLog{
			Ip:        c.ClientIP(),
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			UserAgent: c.Request.UserAgent(),
			RequestID: mws.GetRequestID(c),
			UserID:    userID,
			Status:    c.Writer.Status(),
			Latency:   time.Since(start).Milliseconds(),
			Response:  respBody.String(),
		}
		// 处理请求内容
		if strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
			record.Request = "multipart/form-data"
		} else {
			if len(body) > 1024 {
				record.Request = "[超出记录长度]"
			} else {
				record.Request = string(body)
			}
		}
		// 异步写入数据库，减少请求阻塞
		go func(db *gorm.DB, r models.SysAccessLog) {
			db.Create(&r)
		}(db, record)
	}
}

// bodyWriter 用于捕获响应 body
type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
