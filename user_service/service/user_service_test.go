package service

import (
	"context"
	"gin/proto/generated/user"
	"gin/user_service/model"
	"gin/user_service/repository"
	"strings"
	"testing"

	"gorm.io/gorm"
)

// Mock Repository để test UserService - implement interface đúng
type MockUserRepository struct {
	users map[uint]*model.User
}

// Đảm bảo MockUserRepository implement UserRepositoryInterface
var _ repository.UserRepositoryInterface = (*MockUserRepository)(nil)

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[uint]*model.User),
	}
}

func (m *MockUserRepository) Create(user *model.User) error {
	user.ID = uint(len(m.users) + 1)
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) CreateAccount(account *model.Account) error {
	// Mock implementation - không cần lưu account trong test này
	return nil
}

func (m *MockUserRepository) GetByID(id uint) (*model.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockUserRepository) GetBySDT(sdt string) (*model.User, error) {
	for _, user := range m.users {
		if user.SDT == sdt {
			return user, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockUserRepository) Update(user *model.User) error {
	if _, exists := m.users[user.ID]; exists {
		m.users[user.ID] = user
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (m *MockUserRepository) Delete(id uint) error {
	if _, exists := m.users[id]; exists {
		delete(m.users, id)
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (m *MockUserRepository) GetAll(page, limit int) ([]model.User, int64, error) {
	users := make([]model.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, *user)
	}
	return users, int64(len(users)), nil
}

func (m *MockUserRepository) GetRoleByID(id uint) (*model.Role, error) {
	return &model.Role{Model: gorm.Model{ID: id}, Name: "User"}, nil
}

func (m *MockUserRepository) ListRoles(page, limit int) ([]model.Role, int64, error) {
	roles := []model.Role{
		{Model: gorm.Model{ID: 1}, Name: "Admin"},
		{Model: gorm.Model{ID: 2}, Name: "User"},
	}
	return roles, int64(len(roles)), nil
}

// Unit Test cho CreateUser - Test validation logic THẬT
func TestUserService_CreateUser_Validation(t *testing.T) {
	// Arrange - Tạo service THẬT với mock repository
	mockRepo := NewMockUserRepository()
	userService := NewUserService(mockRepo)

	t.Run("Báo lỗi khi Name rỗng", func(t *testing.T) {
		// Arrange
		req := &user.CreateUserRequest{
			Name:   "", // Name rỗng
			Sdt:    "0123456789",
			RoleId: 1,
		}

		// Act - Gọi hàm THẬT
		response, err := userService.CreateUser(context.Background(), req)

		// Assert - Kiểm tra kết quả THẬT
		if err == nil {
			t.Error("Expected error for empty name, got nil")
		}
		if response != nil {
			t.Error("Expected nil response for empty name, got response")
		}
		if !strings.Contains(err.Error(), "name and SDT are required") {
			t.Errorf("Expected error to contain 'name and SDT are required', got '%s'", err.Error())
		}
	})

	t.Run("Báo lỗi khi SDT rỗng", func(t *testing.T) {
		// Arrange
		req := &user.CreateUserRequest{
			Name:   "John Doe",
			Sdt:    "", // SDT rỗng
			RoleId: 1,
		}

		// Act - Gọi hàm THẬT
		response, err := userService.CreateUser(context.Background(), req)

		// Assert - Kiểm tra kết quả THẬT
		if err == nil {
			t.Error("Expected error for empty SDT, got nil")
		}
		if response != nil {
			t.Error("Expected nil response for empty SDT, got response")
		}
		if !strings.Contains(err.Error(), "name and SDT are required") {
			t.Errorf("Expected error to contain 'name and SDT are required', got '%s'", err.Error())
		}
	})

	t.Run("Báo lỗi khi Name quá dài", func(t *testing.T) {
		// Arrange
		req := &user.CreateUserRequest{
			Name:   "Very very very very long name here", // Name quá dài
			Sdt:    "0123456789",
			RoleId: 1,
		}

		// Act - Gọi hàm THẬT
		response, err := userService.CreateUser(context.Background(), req)

		// Assert - Kiểm tra kết quả THẬT
		if err == nil {
			t.Error("Expected error for long name, got nil")
		}
		if response != nil {
			t.Error("Expected nil response for long name, got response")
		}
		if !strings.Contains(err.Error(), "name must be between 2 and 20 characters") {
			t.Errorf("Expected error to contain 'name must be between 2 and 20 characters', got '%s'", err.Error())
		}
	})

	t.Run("Báo lỗi khi Name quá ngắn", func(t *testing.T) {
		// Arrange
		req := &user.CreateUserRequest{
			Name:   "Jo", // Name quá ngắn
			Sdt:    "0123456789",
			RoleId: 1,
		}

		// Act - Gọi hàm THẬT
		response, err := userService.CreateUser(context.Background(), req)

		// Assert - Kiểm tra kết quả THẬT
		if err == nil {
			t.Error("Expected error for short name, got nil")
		}
		if response != nil {
			t.Error("Expected nil response for short name, got response")
		}
		if !strings.Contains(err.Error(), "name must be between 2 and 20 characters") {
			t.Errorf("Expected error to contain 'name must be between 2 and 20 characters', got '%s'", err.Error())
		}
	})

	t.Run("Báo lỗi khi Name chứa số", func(t *testing.T) {
		// Arrange
		req := &user.CreateUserRequest{
			Name:   "John44", // Name chứa số
			Sdt:    "0123456789",
			RoleId: 1,
		}

		// Act - Gọi hàm THẬT
		response, err := userService.CreateUser(context.Background(), req)

		// Assert - Kiểm tra kết quả THẬT
		// Test này sẽ PASS vì logic check Name chứa số đã được uncomment
		if err == nil {
			t.Error("Expected error for name with numbers, got nil")
		}
		if response != nil {
			t.Error("Expected nil response for name with numbers, got response")
		}
		if !strings.Contains(err.Error(), "name cannot contain numbers") {
			t.Errorf("Expected error to contain 'name cannot contain numbers', got '%s'", err.Error())
		}
		t.Logf("✅ Test PASS: Logic check Name chứa số hoạt động đúng")
	})

	t.Run("Không báo lỗi khi tất cả fields hợp lệ", func(t *testing.T) {
		// Arrange
		req := &user.CreateUserRequest{
			Name:   "John Doe",
			Sdt:    "0123456789",
			RoleId: 1,
		}

		// Act - Gọi hàm THẬT
		response, err := userService.CreateUser(context.Background(), req)

		// Assert - Kiểm tra kết quả THẬT
		if err != nil {
			t.Errorf("Expected no error for valid fields, got error: %v", err)
		}
		if response == nil {
			t.Error("Expected response for valid fields, got nil")
		}
		t.Logf("Validation PASS: Tất cả fields hợp lệ")
	})
}
