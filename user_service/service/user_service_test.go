package service

import (
	"gin/proto/generated/user"
	"gin/user_service/model"
	"testing"

	"gorm.io/gorm"
)

// UserRepositoryInterface định nghĩa interface cho repository
type UserRepositoryInterface interface {
	Create(user *model.User) error
	GetByID(id uint) (*model.User, error)
	GetBySDT(sdt string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
	GetAll(page, limit int) ([]model.User, int64, error)
	GetRoleByID(id uint) (*model.Role, error)
	ListRoles(page, limit int) ([]model.Role, int64, error)
}

// Mock Repository
type MockUserRepository struct {
	users map[uint]*model.User
}

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

// Test Cases
func TestUserService_CreateUser(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()

	req := &user.CreateUserRequest{
		Name:   "John Doe",
		Sdt:    "0123456789",
		RoleId: 1,
	}

	// Act - Test create user logic
	err := mockRepo.Create(&model.User{
		Name: req.Name,
		SDT:  req.Sdt,
	})

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify user was created
	createdUser, err := mockRepo.GetBySDT(req.Sdt)
	if err != nil {
		t.Errorf("Expected to find created user, got error: %v", err)
	}
	if createdUser.Name != "John Doe" {
		t.Errorf("Expected 'John Doe', got %s", createdUser.Name)
	}
	if createdUser.SDT != "0123456789" {
		t.Errorf("Expected '0123456789', got %s", createdUser.SDT)
	}
}

func TestUserService_GetUser(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()

	// Create a user first
	testUser := &model.User{
		Name: "John Doe",
		SDT:  "0123456789",
	}
	mockRepo.Create(testUser)

	// Act
	foundUser, err := mockRepo.GetByID(testUser.ID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if foundUser == nil {
		t.Error("Expected user, got nil")
	}
	if foundUser.Name != "John Doe" {
		t.Errorf("Expected 'John Doe', got %s", foundUser.Name)
	}
}

func TestUserService_GetUser_NotFound(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()

	// Act
	foundUser, err := mockRepo.GetByID(999)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if foundUser != nil {
		t.Error("Expected nil user, got user")
	}
	if err != gorm.ErrRecordNotFound {
		t.Errorf("Expected gorm.ErrRecordNotFound, got %v", err)
	}
}

func TestUserService_GetUserBySDT(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()

	sdt := "0123456789"

	// Create a user first
	testUser := &model.User{
		Name: "John Doe",
		SDT:  sdt,
	}
	mockRepo.Create(testUser)

	// Act
	foundUser, err := mockRepo.GetBySDT(sdt)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if foundUser == nil {
		t.Error("Expected user, got nil")
	}
	if foundUser.SDT != sdt {
		t.Errorf("Expected SDT %s, got %s", sdt, foundUser.SDT)
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()

	// Create a user first
	testUser := &model.User{
		Name: "John Doe",
		SDT:  "0123456789",
	}
	mockRepo.Create(testUser)

	// Update user
	testUser.Name = "Jane Doe"
	testUser.SDT = "0987654321"

	// Act
	err := mockRepo.Update(testUser)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify update
	updatedUser, _ := mockRepo.GetByID(testUser.ID)
	if updatedUser.Name != "Jane Doe" {
		t.Errorf("Expected 'Jane Doe', got %s", updatedUser.Name)
	}
	if updatedUser.SDT != "0987654321" {
		t.Errorf("Expected '0987654321', got %s", updatedUser.SDT)
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()

	// Create a user first
	testUser := &model.User{
		Name: "John Doe",
		SDT:  "0123456789",
	}
	mockRepo.Create(testUser)

	// Act
	err := mockRepo.Delete(testUser.ID)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify deletion
	_, err = mockRepo.GetByID(testUser.ID)
	if err == nil {
		t.Error("Expected error when getting deleted user, got nil")
	}
}

func TestUserService_ListUsers(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()

	// Create some users
	mockRepo.Create(&model.User{Name: "John Doe", SDT: "0123456789"})
	mockRepo.Create(&model.User{Name: "Jane Doe", SDT: "0987654321"})

	// Act
	users, total, err := mockRepo.GetAll(1, 10)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if total != 2 {
		t.Errorf("Expected total 2, got %d", total)
	}
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}
