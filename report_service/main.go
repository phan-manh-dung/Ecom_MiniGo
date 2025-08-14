package main

import (
	"log"
	"net/http"
	"os"

	"report-service/db"
	"report-service/handler"
	"report-service/repository"
	"report-service/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Connect to database
	db.ConnectDatabase()

	// Initialize repository
	reportRepo := repository.NewReportRepository()

	// Initialize services
	reportService := service.NewReportService(reportRepo)

	// Initialize handlers
	reportHandler := handler.NewReportHandler(reportService)

	// Create Gin router
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Setup routes
	setupRoutes(r, reportHandler)

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "30051"
	}

	log.Printf("Report Service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRoutes(r *gin.Engine, reportHandler *handler.ReportHandler) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "report-service"})
	})

	// API routes
	api := r.Group("/api")
	{
		// Admin reports (cần authentication)
		admin := api.Group("/admin")
		{
			reports := admin.Group("/reports")
			{
				reports.POST("/generate", reportHandler.GenerateReport)
			}
		}
	}
}
