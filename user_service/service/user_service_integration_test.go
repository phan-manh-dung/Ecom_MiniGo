package service

import (
	"context"
	"fmt"
	"gin/user_service/model"
	"gin/user_service/repository"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// UserServiceIntegrationTestSuite test suite cho integration testing với database thật
type UserServiceIntegrationTestSuite struct {
	suite.Suite
	db       *gorm.DB
	userRepo *repository.UserRepository
	ctx      context.Context
}

// SetupSuite chạy một lần trước tất cả tests
func (suite *UserServiceIntegrationTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	// Lấy database URL từ environment hoặc dùng default test database
	databaseURL := getTestDatabaseURL()

	// Kết nối database thật (PostgreSQL)
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Tắt log trong test
	})
	if err != nil {
		suite.T().Fatalf("Failed to connect to test database: %v", err)
	}

	suite.db = db

	// Auto migrate tables theo thứ tự để tránh foreign key constraint issues
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
	// Cleanup data cũ trước mỗi test để tránh conflict
	err := suite.cleanupTestData()
	if err != nil {
		suite.T().Fatalf("Failed to cleanup test data: %v", err)
	}
}

// TearDownSuite chạy một lần sau tất cả tests
func (suite *UserServiceIntegrationTestSuite) TearDownSuite() {
	// Cleanup cuối cùng
	err := suite.cleanupTestData()
	if err != nil {
		suite.T().Logf("Warning: Failed to cleanup test data: %v", err)
	}
}

// cleanupTestData xóa tất cả test data để tránh conflict
func (suite *UserServiceIntegrationTestSuite) cleanupTestData() error {
	// Disable foreign key checks tạm thời
	err := suite.db.Exec("SET session_replication_role = replica").Error
	if err != nil {
		return err
	}

	// Truncate tables theo thứ tự để tránh foreign key constraint
	err = suite.db.Exec("TRUNCATE TABLE accounts RESTART IDENTITY CASCADE").Error
	if err != nil {
		return err
	}

	err = suite.db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE").Error
	if err != nil {
		return err
	}

	err = suite.db.Exec("TRUNCATE TABLE roles RESTART IDENTITY CASCADE").Error
	if err != nil {
		return err
	}

	// Re-enable foreign key checks
	err = suite.db.Exec("SET session_replication_role = DEFAULT").Error
	if err != nil {
		return err
	}

	return nil
}

// generateUniqueSDT tạo SDT duy nhất cho mỗi test
func (suite *UserServiceIntegrationTestSuite) generateUniqueSDT() string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("test_%d", timestamp)
}

// generateUniqueName tạo tên duy nhất cho mỗi test
func (suite *UserServiceIntegrationTestSuite) generateUniqueName() string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("Test User %d", timestamp)
}

// TestCreateUser test tạo user mới
func (suite *UserServiceIntegrationTestSuite) TestCreateUser() {
	// Arrange - Tạo user test với data duy nhất
	testUser := &model.User{
		Name: suite.generateUniqueName(),
		SDT:  suite.generateUniqueSDT(),
	}

	// Act - Tạo user qua repository
	err := suite.userRepo.Create(testUser)

	// Assert - Kiểm tra không có lỗi
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), testUser.ID) // ID phải được tạo

	// Verify - Kiểm tra user đã được tạo
	foundUser, err := suite.userRepo.GetBySDT(testUser.SDT)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), testUser.Name, foundUser.Name)
	assert.Equal(suite.T(), testUser.SDT, foundUser.SDT)
}

// TestGetUserBySDT test tìm user theo SDT
func (suite *UserServiceIntegrationTestSuite) TestGetUserBySDT() {
	// Arrange - Tạo user test
	testUser := &model.User{
		Name: suite.generateUniqueName(),
		SDT:  suite.generateUniqueSDT(),
	}

	err := suite.userRepo.Create(testUser)
	assert.NoError(suite.T(), err)

	// Act - Tìm user theo SDT
	foundUser, err := suite.userRepo.GetBySDT(testUser.SDT)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), testUser.Name, foundUser.Name)
	assert.Equal(suite.T(), testUser.SDT, foundUser.SDT)
}

// TestGetUserBySDT_NotFound test trường hợp không tìm thấy user
func (suite *UserServiceIntegrationTestSuite) TestGetUserBySDT_NotFound() {
	// Act - Tìm user không tồn tại
	foundUser, err := suite.userRepo.GetBySDT("non_existent_sdt")

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), foundUser)
}

// TestUpdateUser test cập nhật user
func (suite *UserServiceIntegrationTestSuite) TestUpdateUser() {
	// Arrange - Tạo user test
	testUser := &model.User{
		Name: suite.generateUniqueName(),
		SDT:  suite.generateUniqueSDT(),
	}

	err := suite.userRepo.Create(testUser)
	assert.NoError(suite.T(), err)

	// Act - Cập nhật user
	newName := suite.generateUniqueName()
	testUser.Name = newName
	err = suite.userRepo.Update(testUser)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify - Kiểm tra user đã được cập nhật
	foundUser, err := suite.userRepo.GetByID(testUser.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), newName, foundUser.Name)
}

// TestDeleteUser test xóa user
func (suite *UserServiceIntegrationTestSuite) TestDeleteUser() {
	// Arrange - Tạo user test
	testUser := &model.User{
		Name: suite.generateUniqueName(),
		SDT:  suite.generateUniqueSDT(),
	}

	err := suite.userRepo.Create(testUser)
	assert.NoError(suite.T(), err)
	userID := testUser.ID

	// Act - Xóa user
	err = suite.userRepo.Delete(userID)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify - Kiểm tra user đã bị xóa
	foundUser, err := suite.userRepo.GetByID(userID)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), foundUser)
}

// TestListUsers test lấy danh sách users
func (suite *UserServiceIntegrationTestSuite) TestListUsers() {
	// Arrange - Tạo nhiều users test
	for i := 0; i < 3; i++ {
		testUser := &model.User{
			Name: suite.generateUniqueName(),
			SDT:  suite.generateUniqueSDT(),
		}
		err := suite.userRepo.Create(testUser)
		assert.NoError(suite.T(), err)
	}

	// Act - Lấy danh sách users
	users, total, err := suite.userRepo.GetAll(1, 10)

	// Assert
	assert.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), total, int64(3))
	assert.GreaterOrEqual(suite.T(), len(users), 3)
}

// getTestDatabaseURL lấy database URL cho test
func getTestDatabaseURL() string {
	// Ưu tiên environment variable
	if url := os.Getenv("TEST_DATABASE_URL"); url != "" {
		return url
	}

	// Default test database - KHÔNG phải production
	return "host=localhost user=postgres password=123postgres dbname=user_service_test port=5432 sslmode=disable"
}

// Chạy test suite
func TestUserServiceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceIntegrationTestSuite))
}
