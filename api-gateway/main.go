package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Sử dụng Wire để khởi tạo application với dependency injection
	app, err := wireApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

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
