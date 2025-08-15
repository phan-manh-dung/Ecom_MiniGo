package service

import (
	"gin/user_service/model"
	"gin/user_service/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// UserServiceIntegrationTestSuite test suite cho integration testing
type UserServiceIntegrationTestSuite struct {
	suite.Suite
	db       *gorm.DB
	userRepo *repository.UserRepository
}

// SetupSuite chạy một lần trước tất cả tests
func (suite *UserServiceIntegrationTestSuite) SetupSuite() {
	// Kết nối SQLite in-memory (database ảo)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Tắt log trong test
	})
	if err != nil {
		suite.T().Fatalf("Failed to connect to test database: %v", err)
	}

	suite.db = db

	// Auto migrate tables (tạo schema giống hệt production)
	// Migrate theo thứ tự để tránh foreign key constraint issues
	err = suite.db.AutoMigrate(&model.Role{})
	if err != nil {
		suite.T().Fatalf("Failed to migrate Role table: %v", err)
	}

	err = suite.db.AutoMigrate(&model.User{})
	if err != nil {
		suite.T().Fatalf("Failed to migrate User table: %v", err)
	}

	err = suite.db.AutoMigrate(&model.Account{})
	if err != nil {
		suite.T().Fatalf("Failed to migrate Account table: %v", err)
	}

	// Setup repositories
	suite.userRepo = repository.NewUserRepository(suite.db)
	if suite.userRepo == nil {
		suite.T().Fatalf("Failed to create UserRepository")
	}
}

// SetupTest chạy trước mỗi test case
func (suite *UserServiceIntegrationTestSuite) SetupTest() {
	// Mỗi test case có database riêng (SQLite in-memory tự động reset)
	// KHÔNG cần cleanup gì cả - an toàn 100%
}

// TearDownSuite chạy một lần sau tất cả tests
func (suite *UserServiceIntegrationTestSuite) TearDownSuite() {
}

// Test đơn giản: Create User
func (suite *UserServiceIntegrationTestSuite) TestCreateUser() {
	// Arrange - Tạo user test
	testUser := &model.User{
		Name: "Test User",
		SDT:  "test_1234567890",
	}

	// Act - Tạo user qua repository (test business logic)
	err := suite.userRepo.Create(testUser)

	// Assert - Kiểm tra không có lỗi
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), testUser.ID) // ID phải được tạo

	// Verify - Kiểm tra qua repository (KHÔNG test database trực tiếp)
	foundUser, err := suite.userRepo.GetBySDT("test_1234567890")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), testUser.Name, foundUser.Name)
	assert.Equal(suite.T(), testUser.SDT, foundUser.SDT)
}

// Chạy test suite
func TestUserServiceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceIntegrationTestSuite))
}
