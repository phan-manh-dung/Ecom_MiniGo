package service

import (
	"gin/user_service/model"
	"gin/user_service/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
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
	// Kết nối test database
	dsn := "host=localhost user=postgres password=123postgres dbname=user_service port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Tắt log trong test
	})
	if err != nil {
		suite.T().Fatalf("Failed to connect to test database: %v", err)
	}

	suite.db = db

	// Auto migrate tables
	err = suite.db.AutoMigrate(&model.User{}, &model.Role{})
	if err != nil {
		suite.T().Fatalf("Failed to migrate test database: %v", err)
	}

	// Setup repositories
	suite.userRepo = repository.NewUserRepository(suite.db)
}

// SetupTest chạy trước mỗi test case
func (suite *UserServiceIntegrationTestSuite) SetupTest() {
	// KHÔNG xóa data thật! Chỉ xóa test data cụ thể
	suite.db.Exec("DELETE FROM users WHERE sdt LIKE 'test_%' OR sdt = '1234567890'")
}

// TearDownSuite chạy một lần sau tất cả tests
func (suite *UserServiceIntegrationTestSuite) TearDownSuite() {
	// KHÔNG xóa gì cả! Giữ nguyên data thật
	suite.db.Exec("DELETE FROM users WHERE sdt LIKE 'test_%' OR sdt = '1234567890'")
}

// Test đơn giản: Create User
func (suite *UserServiceIntegrationTestSuite) TestCreateUser() {
	// Arrange - Tạo user test với SDT an toàn
	testUser := &model.User{
		Name: "Test User",
		SDT:  "test_1234567890", // SDT bắt đầu bằng "test_" để dễ nhận biết
	}

	// Act - Tạo user qua repository
	err := suite.userRepo.Create(testUser)

	// Assert - Kiểm tra không có lỗi
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), testUser.ID) // ID phải được tạo

	// Verify - Kiểm tra user đã được lưu vào database
	var createdUser model.User
	err = suite.db.Where("sdt = ?", testUser.SDT).First(&createdUser).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), testUser.Name, createdUser.Name)
	assert.Equal(suite.T(), testUser.SDT, createdUser.SDT)
}

// Chạy test suite
func TestUserServiceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceIntegrationTestSuite))
}
