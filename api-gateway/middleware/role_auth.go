package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
