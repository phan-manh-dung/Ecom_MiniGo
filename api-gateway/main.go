package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// wireApp sẽ được generate bởi Wire tool
// Chạy: wire
func main() {
	// Sử dụng Wire để khởi tạo application với dependency injection
	app, err := wireApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	//  test
	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "API Gateway is running",
		})
	})

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway starting on port %s", port)
	if err := app.Run(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
