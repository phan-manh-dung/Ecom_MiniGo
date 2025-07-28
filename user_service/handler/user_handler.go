package handler

import (
	"context"

	"gin/proto/generated/user"
	"gin/user_service/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	userService *service.UserService
	user.UnimplementedUserServiceServer
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUser implements user.UserServiceServer
func (h *UserHandler) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	response, err := h.userService.GetUser(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}
	return response, nil
}

// GetUserBySDT implements user.UserServiceServer
func (h *UserHandler) GetUserBySDT(ctx context.Context, req *user.GetUserBySDTRequest) (*user.GetUserBySDTResponse, error) {
	resp, err := h.userService.GetUserBySDT(ctx, req.Sdt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user by SDT: %v", err)
	}
	return resp, nil
}

// CreateUser implements user.UserServiceServer
func (h *UserHandler) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	resp, err := h.userService.CreateUser(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}
	return resp, nil
}

// UpdateUser implements user.UserServiceServer
func (h *UserHandler) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	resp, err := h.userService.UpdateUser(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}
	return resp, nil
}

// DeleteUser implements user.UserServiceServer
func (h *UserHandler) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
	resp, err := h.userService.DeleteUser(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}
	return resp, nil
}

// ListUsers implements user.UserServiceServer
func (h *UserHandler) ListUsers(ctx context.Context, req *user.ListUsersRequest) (*user.ListUsersResponse, error) {
	resp, err := h.userService.ListUsers(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}
	return resp, nil
}

// GetRole implements user.UserServiceServer
func (h *UserHandler) GetRole(ctx context.Context, req *user.GetRoleRequest) (*user.GetRoleResponse, error) {
	resp, err := h.userService.GetRole(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get role: %v", err)
	}
	return resp, nil
}

// ListRoles implements user.UserServiceServer
func (h *UserHandler) ListRoles(ctx context.Context, req *user.ListRolesRequest) (*user.ListRolesResponse, error) {
	resp, err := h.userService.ListRoles(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list roles: %v", err)
	}
	return resp, nil
}
