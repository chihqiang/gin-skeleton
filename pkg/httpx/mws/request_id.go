package mws

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	headerXRequestID = "X-Request-ID"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(headerXRequestID)
		if rid == "" {
			rid = uuid.New().String()
			c.Request.Header.Add(headerXRequestID, rid)
		}
		c.Header(headerXRequestID, rid)
		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	return c.GetHeader(headerXRequestID)
}
