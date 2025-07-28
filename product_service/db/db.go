package db

import (
	"fmt"
	"log"
	"os"

	models "gin/product_service/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Thông tin kết nối database
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "123postgres")
	dbname := getEnv("DB_NAME", "product_service")
	sslmode := getEnv("DB_SSLMODE", "disable")

	// Tạo DSN (Data Source Name) cho PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Ho_Chi_Minh",
		host, user, password, dbname, port, sslmode)

	// Cấu hình GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Log SQL queries
	}

	// Kết nối database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate các model
	err = DB.AutoMigrate(
		&models.Product{},
		&models.Inventory{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected successfully!")
}

// GetDB trả về instance của database
func GetDB() *gorm.DB {
	return DB
}

// getEnv lấy giá trị từ environment variable, nếu không có thì dùng default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
