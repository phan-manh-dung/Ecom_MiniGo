package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// NewRequestIDMiddleware tạo Request ID middleware
func NewRequestIDMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Tạo request ID nếu chưa có
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		// Set request ID vào header
		c.Header("X-Request-ID", requestID)
		
		// Set vào context để các handler có thể sử dụng
		c.Set("request_id", requestID)
		
		c.Next()
	})
}

