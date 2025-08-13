package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// NewAuthMiddleware tạo Authentication middleware
func NewAuthMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Bỏ qua auth cho một số routes
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/api/auth/login" {
			c.Next()
			return
		}
		
		// Lấy token từ header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}
		
		// Kiểm tra format: "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format. Use: Bearer <token>",
			})
			c.Abort()
			return
		}
		
		token := tokenParts[1]
		
		// TODO: Validate JWT token
		// Đây là placeholder, bạn cần implement JWT validation
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}
		
		// Token hợp lệ, tiếp tục
		c.Next()
	})
}

// phần check role ADMIN
// AdminOnlyMiddleware chỉ cho phép ADMIN access
func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy role từ context (đã được set bởi AuthMiddleware)
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role not found in context"})
			c.Abort()
			return
		}

		// Kiểm tra role
		if role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Admin role required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RoleMiddleware kiểm tra role cụ thể
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role not found in context"})
			c.Abort()
			return
		}

		// Kiểm tra role có trong danh sách được phép không
		roleStr := role.(string)
		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Insufficient permissions"})
		c.Abort()
	}
}
