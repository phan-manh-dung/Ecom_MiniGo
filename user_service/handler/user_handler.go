package handler

// Implement các method từ protobuf definition
// Protocol Translation: Chuyển đổi giữa gRPC request/response và service layer
// Error Handling: Xử lý lỗi và trả về gRPC status codes chuẩn
// Request Delegation: Chuyển request xuống service layer xử lý
import (
	"context"

	"gin/proto/generated/user"
	"gin/shared/generic"
	"gin/user_service/service"
)

type UserHandler struct {
	userService *service.UserService
	generic     *generic.GenericHandler
	user.UnimplementedUserServiceServer
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		generic:     generic.NewGenericHandler(),
	}
}

// GetUser implements user.UserServiceServer
func (h *UserHandler) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	return generic.HandleOperationWithID[uint32, *user.GetUserResponse, uint32](ctx, req.Id, h.userService.GetUser, "get user")
}

// GetUserBySDT implements user.UserServiceServer
func (h *UserHandler) GetUserBySDT(ctx context.Context, req *user.GetUserBySDTRequest) (*user.GetUserBySDTResponse, error) {
	return generic.HandleOperationWithString[*user.GetUserBySDTRequest, *user.GetUserBySDTResponse](ctx, req.Sdt, h.userService.GetUserBySDT, "get user by SDT")
}

// CreateUser implements user.UserServiceServer
func (h *UserHandler) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	return generic.HandleOperation[*user.CreateUserRequest, *user.CreateUserResponse](ctx, req, h.userService.CreateUser, "create user")
}

// UpdateUser implements user.UserServiceServer
func (h *UserHandler) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	return generic.HandleOperation[*user.UpdateUserRequest, *user.UpdateUserResponse](ctx, req, h.userService.UpdateUser, "update user")
}

// DeleteUser implements user.UserServiceServer
func (h *UserHandler) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
	return generic.HandleOperation[*user.DeleteUserRequest, *user.DeleteUserResponse](ctx, req, h.userService.DeleteUser, "delete user")
}

// ListUsers implements user.UserServiceServer
func (h *UserHandler) ListUsers(ctx context.Context, req *user.ListUsersRequest) (*user.ListUsersResponse, error) {
	return generic.HandleListOperation[*user.ListUsersRequest, *user.ListUsersResponse](ctx, req, h.userService.ListUsers, "list users")
}

// GetRole implements user.UserServiceServer
func (h *UserHandler) GetRole(ctx context.Context, req *user.GetRoleRequest) (*user.GetRoleResponse, error) {
	return generic.HandleOperation[*user.GetRoleRequest, *user.GetRoleResponse](ctx, req, h.userService.GetRole, "get role")
}

// ListRoles implements user.UserServiceServer
func (h *UserHandler) ListRoles(ctx context.Context, req *user.ListRolesRequest) (*user.ListRolesResponse, error) {
	return generic.HandleListOperation[*user.ListRolesRequest, *user.ListRolesResponse](ctx, req, h.userService.ListRoles, "list roles")
}
