package service

import (
	"gin/user_service/model"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// TestDB interface để test database operations
type TestDB interface {
	Create(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	Exec(sql string, values ...interface{}) *gorm.DB
}

// TestDataManager quản lý test data
type TestDataManager struct {
	db TestDB
}

func NewTestDataManager(db TestDB) *TestDataManager {
	return &TestDataManager{db: db}
}

// CreateTestUser tạo user test với data mặc định
func (tdm *TestDataManager) CreateTestUser(t *testing.T) *model.User {
	user := &model.User{
		Name: "Test User",
		SDT:  "0123456789",
	}

	err := tdm.db.Create(user).Error
	require.NoError(t, err)

	return user
}

// CreateTestUserWithData tạo user test với data tùy chỉnh
func (tdm *TestDataManager) CreateTestUserWithData(t *testing.T, name, sdt string) *model.User {
	user := &model.User{
		Name: name,
		SDT:  sdt,
	}

	err := tdm.db.Create(user).Error
	require.NoError(t, err)

	return user
}

// CreateTestRole tạo role test
func (tdm *TestDataManager) CreateTestRole(t *testing.T, name string) *model.Role {
	role := &model.Role{Name: name}

	err := tdm.db.Create(role).Error
	require.NoError(t, err)

	return role
}

// CleanupTestData xóa tất cả test data
func (tdm *TestDataManager) CleanupTestData(t *testing.T) {
	err := tdm.db.Exec("DELETE FROM users").Error
	require.NoError(t, err)

	err = tdm.db.Exec("DELETE FROM roles").Error
	require.NoError(t, err)
}

// AssertUserExists kiểm tra user tồn tại trong database
func (tdm *TestDataManager) AssertUserExists(t *testing.T, sdt string) *model.User {
	var user model.User
	err := tdm.db.Where("sdt = ?", sdt).First(&user).Error
	require.NoError(t, err)
	return &user
}

// AssertUserNotExists kiểm tra user không tồn tại
func (tdm *TestDataManager) AssertUserNotExists(t *testing.T, sdt string) {
	var user model.User
	err := tdm.db.Where("sdt = ?", sdt).First(&user).Error
	require.Error(t, err)
	require.Equal(t, gorm.ErrRecordNotFound, err)
}

// TestScenario struct để định nghĩa test scenarios
type TestScenario struct {
	Name        string
	Description string
	Setup       func(t *testing.T) interface{}
	Action      func(t *testing.T, data interface{}) error
	Assert      func(t *testing.T, data interface{}, err error)
}

// RunTestScenarios chạy nhiều test scenarios
func RunTestScenarios(t *testing.T, scenarios []TestScenario) {
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			// Setup
			data := scenario.Setup(t)

			// Action
			err := scenario.Action(t, data)

			// Assert
			scenario.Assert(t, data, err)
		})
	}
}

// TestDataFactory factory để tạo test data
type TestDataFactory struct{}

// NewTestDataFactory tạo instance mới
func NewTestDataFactory() *TestDataFactory {
	return &TestDataFactory{}
}

// GenerateUsers tạo nhiều users test
func (f *TestDataFactory) GenerateUsers(count int) []*model.User {
	users := make([]*model.User, count)
	for i := 0; i < count; i++ {
		users[i] = &model.User{
			Name: "User " + string(rune('A'+i)),
			SDT:  "012345678" + string(rune('0'+i)),
		}
	}
	return users
}

// GenerateSDTs tạo danh sách SDT test
func (f *TestDataFactory) GenerateSDTs(count int) []string {
	sdts := make([]string, count)
	for i := 0; i < count; i++ {
		sdts[i] = "012345678" + string(rune('0'+i))
	}
	return sdts
}

// TestValidationHelper helper để test validation
type TestValidationHelper struct{}

// NewTestValidationHelper tạo instance mới
func NewTestValidationHelper() *TestValidationHelper {
	return &TestValidationHelper{}
}

// AssertValidationError kiểm tra validation error
func (h *TestValidationHelper) AssertValidationError(t *testing.T, err error, expectedField string) {
	require.Error(t, err)
	require.Contains(t, err.Error(), expectedField)
}

// AssertNoValidationError kiểm tra không có validation error
func (h *TestValidationHelper) AssertNoValidationError(t *testing.T, err error) {
	require.NoError(t, err)
}
