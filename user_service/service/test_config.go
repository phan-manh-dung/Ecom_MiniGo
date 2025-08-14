package service

import (
	"gin/user_service/model"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestConfig cấu hình cho test environment
type TestConfig struct {
	DatabaseURL string
	LogLevel    logger.LogLevel
	AutoMigrate bool
}

// DefaultTestConfig cấu hình mặc định cho test
func DefaultTestConfig() *TestConfig {
	return &TestConfig{
		DatabaseURL: getTestDatabaseURL(),
		LogLevel:    logger.Silent,
		AutoMigrate: true,
	}
}

// getTestDatabaseURL lấy database URL từ environment hoặc dùng default
func getTestDatabaseURL() string {
	if url := os.Getenv("TEST_DATABASE_URL"); url != "" {
		return url
	}

	// Default test database
	return "host=localhost user=postgres password=123postgres dbname=user_service port=5432 sslmode=disable"
}

// TestDatabase test database instance
type TestDatabase struct {
	DB     *gorm.DB
	Config *TestConfig
}

// NewTestDatabase tạo test database mới
func NewTestDatabase(config *TestConfig) (*TestDatabase, error) {
	if config == nil {
		config = DefaultTestConfig()
	}

	db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
	})
	if err != nil {
		return nil, err
	}

	testDB := &TestDatabase{
		DB:     db,
		Config: config,
	}

	// Auto migrate nếu được yêu cầu
	if config.AutoMigrate {
		err = testDB.Migrate()
		if err != nil {
			return nil, err
		}
	}

	return testDB, nil
}

// Migrate thực hiện database migration
func (tdb *TestDatabase) Migrate() error {
	return tdb.DB.AutoMigrate(
		&model.User{},
		&model.Role{},
	)
}

// Cleanup xóa tất cả data trong test database
func (tdb *TestDatabase) Cleanup() error {
	// Disable foreign key checks
	err := tdb.DB.Exec("SET session_replication_role = replica").Error
	if err != nil {
		return err
	}

	// Truncate tables
	err = tdb.DB.Exec("TRUNCATE TABLE users, roles RESTART IDENTITY CASCADE").Error
	if err != nil {
		return err
	}

	// Re-enable foreign key checks
	err = tdb.DB.Exec("SET session_replication_role = DEFAULT").Error
	if err != nil {
		return err
	}

	return nil
}

// Close đóng database connection
func (tdb *TestDatabase) Close() error {
	sqlDB, err := tdb.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// TestSuiteBase base struct cho tất cả test suites
type TestSuiteBase struct {
	DB               *TestDatabase
	DataManager      *TestDataManager
	DataFactory      *TestDataFactory
	ValidationHelper *TestValidationHelper
}

// SetupTestSuite setup cho test suite
func (tsb *TestSuiteBase) SetupTestSuite(t *testing.T) {
	// Tạo test database
	testDB, err := NewTestDatabase(nil)
	require.NoError(t, err)
	tsb.DB = testDB

	// Tạo test data manager
	tsb.DataManager = NewTestDataManager(testDB.DB)

	// Tạo test data factory
	tsb.DataFactory = NewTestDataFactory()

	// Tạo validation helper
	tsb.ValidationHelper = NewTestValidationHelper()
}

// SetupTest setup cho mỗi test case
func (tsb *TestSuiteBase) SetupTest(t *testing.T) {
	// Cleanup database trước mỗi test
	err := tsb.DB.Cleanup()
	require.NoError(t, err)

	// Tạo test roles
	tsb.DataManager.CreateTestRole(t, "Admin")
	tsb.DataManager.CreateTestRole(t, "User")
}

// TearDownTestSuite cleanup sau test suite
func (tsb *TestSuiteBase) TearDownTestSuite(t *testing.T) {
	if tsb.DB != nil {
		err := tsb.DB.Close()
		require.NoError(t, err)
	}
}

// TestEnvironment test environment configuration
type TestEnvironment struct {
	Database *TestDatabase
	Config   *TestConfig
}

// NewTestEnvironment tạo test environment mới
func NewTestEnvironment(config *TestConfig) (*TestEnvironment, error) {
	testDB, err := NewTestDatabase(config)
	if err != nil {
		return nil, err
	}

	return &TestEnvironment{
		Database: testDB,
		Config:   config,
	}, nil
}

// SetupTestEnvironment setup test environment
func SetupTestEnvironment(t *testing.T) *TestEnvironment {
	env, err := NewTestEnvironment(DefaultTestConfig())
	require.NoError(t, err)

	t.Cleanup(func() {
		err := env.Database.Close()
		require.NoError(t, err)
	})

	return env
}

// TestDatabaseFixture fixture cho test database
type TestDatabaseFixture struct {
	Users []*model.User
	Roles []*model.Role
}

// NewTestDatabaseFixture tạo fixture mới
func NewTestDatabaseFixture() *TestDatabaseFixture {
	return &TestDatabaseFixture{
		Users: make([]*model.User, 0),
		Roles: make([]*model.Role, 0),
	}
}

// AddUser thêm user vào fixture
func (f *TestDatabaseFixture) AddUser(user *model.User) *TestDatabaseFixture {
	f.Users = append(f.Users, user)
	return f
}

// AddRole thêm role vào fixture
func (f *TestDatabaseFixture) AddRole(role *model.Role) *TestDatabaseFixture {
	f.Roles = append(f.Roles, role)
	return f
}

// LoadIntoDatabase load fixture vào database
func (f *TestDatabaseFixture) LoadIntoDatabase(db *gorm.DB) error {
	// Load roles trước
	for _, role := range f.Roles {
		err := db.Create(role).Error
		if err != nil {
			return err
		}
	}

	// Load users sau
	for _, user := range f.Users {
		err := db.Create(user).Error
		if err != nil {
			return err
		}
	}

	return nil
}
