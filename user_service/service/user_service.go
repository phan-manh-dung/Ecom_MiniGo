package service

import (
	"context"
	"fmt"

	"gin/proto/generated/user"
	"gin/shared/generic"
	"gin/user_service/model"
	"gin/user_service/repository"

	"gorm.io/gorm"
)

type UserService struct {
	*generic.BaseService
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		BaseService: generic.NewBaseService(),
		userRepo:    userRepo,
	}
}

// GetUser business logic
func (s *UserService) GetUser(ctx context.Context, id uint32) (*user.GetUserResponse, error) {
	userModel, err := s.userRepo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &user.GetUserResponse{
				User:    nil,
				Message: "User not found",
			}, nil
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	// Convert model to proto
	protoUser := s.convertToProtoUser(userModel)

	return &user.GetUserResponse{
		User:    protoUser,
		Message: "User retrieved successfully",
	}, nil
}

// GetUserBySDT business logic
func (s *UserService) GetUserBySDT(ctx context.Context, sdt string) (*user.GetUserBySDTResponse, error) {
	userModel, err := s.userRepo.GetBySDT(sdt)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &user.GetUserBySDTResponse{
				User:    nil,
				Message: "User not found",
			}, nil
		}
		return nil, fmt.Errorf("failed to get user by SDT: %v", err)
	}

	protoUser := s.convertToProtoUser(userModel)

	return &user.GetUserBySDTResponse{
		User:    protoUser,
		Message: "User retrieved successfully",
	}, nil
}

// CreateUser business logic
func (s *UserService) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	// Validation
	if req.Name == "" || req.Sdt == "" {
		return nil, fmt.Errorf("name and SDT are required")
	}

	// Check if user with SDT already exists
	existingUser, _ := s.userRepo.GetBySDT(req.Sdt)
	if existingUser != nil {
		return nil, fmt.Errorf("user with SDT %s already exists", req.Sdt)
	}

	// Create user model save to database
	userModel := &model.User{
		Name: req.Name,
		SDT:  req.Sdt,
	}

	// Create account if role_id provided
	if req.RoleId > 0 {
		userModel.Account = &model.Account{
			RoleID: uint(req.RoleId),
		}
	}

	// Save to database
	if err := s.userRepo.Create(userModel); err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	// Convert to proto
	protoUser := s.convertToProtoUser(userModel)

	return &user.CreateUserResponse{
		User:    protoUser,
		Message: "User created successfully",
	}, nil
}

// UpdateUser business logic
func (s *UserService) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	userModel, err := s.userRepo.GetByID(uint(req.Id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &user.UpdateUserResponse{
				User:    nil,
				Message: "User not found",
			}, nil
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	// Update fields if provided
	if req.Name != "" {
		userModel.Name = req.Name
	}
	if req.Sdt != "" {
		userModel.SDT = req.Sdt
	}

	if err := s.userRepo.Update(userModel); err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	protoUser := s.convertToProtoUser(userModel)
	return &user.UpdateUserResponse{
		User:    protoUser,
		Message: "User updated successfully",
	}, nil
}

// DeleteUser business logic
func (s *UserService) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
	err := s.userRepo.Delete(uint(req.Id))
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %v", err)
	}
	return &user.DeleteUserResponse{
		Message: "User deleted successfully",
	}, nil
}

// ListUsers business logic
func (s *UserService) ListUsers(ctx context.Context, req *user.ListUsersRequest) (*user.ListUsersResponse, error) {
	// Set default pagination
	page, limit := int(req.Page), int(req.Limit)
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	users, total, err := s.userRepo.GetAll(page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %v", err)
	}

	protoUsers := make([]*user.User, 0, len(users))
	for _, u := range users {
		protoUsers = append(protoUsers, s.convertToProtoUser(&u))
	}

	return &user.ListUsersResponse{
		Users:   protoUsers,
		Total:   int32(total),
		Message: "Users listed successfully",
	}, nil
}

// GetRole business logic
func (s *UserService) GetRole(ctx context.Context, req *user.GetRoleRequest) (*user.GetRoleResponse, error) {
	roleModel, err := s.userRepo.GetRoleByID(uint(req.Id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &user.GetRoleResponse{
				Role:    nil,
				Message: "Role not found",
			}, nil
		}
		return nil, fmt.Errorf("failed to get role: %v", err)
	}

	protoRole := &user.Role{
		Id:   uint32(roleModel.ID),
		Name: roleModel.Name,
	}

	return &user.GetRoleResponse{
		Role:    protoRole,
		Message: "Role retrieved successfully",
	}, nil
}

// ListRoles business logic
func (s *UserService) ListRoles(ctx context.Context, req *user.ListRolesRequest) (*user.ListRolesResponse, error) {
	// Set default pagination
	page, limit := int(req.Page), int(req.Limit)
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	roles, total, err := s.userRepo.ListRoles(page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %v", err)
	}

	protoRoles := make([]*user.Role, 0, len(roles))
	for _, r := range roles {
		protoRoles = append(protoRoles, &user.Role{
			Id:   uint32(r.ID),
			Name: r.Name,
		})
	}

	return &user.ListRolesResponse{
		Roles:   protoRoles,
		Total:   int32(total),
		Message: "Roles listed successfully",
	}, nil
}

// Helper function to convert model.User to proto.User
func (s *UserService) convertToProtoUser(userModel *model.User) *user.User {
	protoUser := &user.User{
		Id:   uint32(userModel.ID),
		Name: userModel.Name,
		Sdt:  userModel.SDT,
	}

	// Add Account info if exists
	if userModel.Account != nil {
		protoUser.Account = &user.Account{
			Id:     uint32(userModel.Account.ID),
			UserId: uint32(userModel.Account.UserID),
			RoleId: uint32(userModel.Account.RoleID),
		}

		// Add Role info if exists
		if userModel.Account.Role.ID != 0 {
			protoUser.Account.Role = &user.Role{
				Id:   uint32(userModel.Account.Role.ID),
				Name: userModel.Account.Role.Name,
			}
		}
	}

	return protoUser
}
