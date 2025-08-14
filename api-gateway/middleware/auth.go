package middleware

import (
	"net/http"
	"strings"

	"gin/api-gateway/utils"

	"github.com/gin-gonic/gin"
)

// NewAuthMiddleware tạo Authentication middleware
func NewAuthMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Bỏ qua auth cho public routes
		publicPaths := []string{
			"/health",
			"/api/auth/login",
			"/api/auth/register", // Nếu có register endpoint
		}
		
		for _, path := range publicPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
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
		
		// Validate JWT token
		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token: " + err.Error(),
			})
			c.Abort()
			return
		}
		
		// Set user info vào context để các middleware khác sử dụng
		c.Set("user_id", claims.UserID)
		c.Set("user_name", claims.Username)
		c.Set("role", claims.Role)
		
		// Token hợp lệ, tiếp tục
		c.Next()
	})
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

// RequireRoleMiddleware kiểm tra role cho route cụ thể
func RequireRoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy role từ context (đã được set bởi AuthMiddleware)
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role not found in context"})
			c.Abort()
			return
		}

		// Kiểm tra role
		if role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. " + requiredRole + " role required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
