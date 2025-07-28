package handler

import (
	"context"
	"fmt"
	"log"

	"gin/proto/generated/user"
	"gin/user_service/model"

	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
	user.UnimplementedUserServiceServer
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

// GetUser implements user.UserServiceServer
func (h *UserHandler) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	log.Printf("GetUser called with ID: %d", req.Id)

	var userModel model.User
	result := h.db.Preload("Account.Role").First(&userModel, req.Id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &user.GetUserResponse{
				User:    nil,
				Message: "User not found",
			}, nil
		}
		return nil, fmt.Errorf("failed to get user: %v", result.Error)
	}

	// Convert model.User to proto.User
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

	return &user.GetUserResponse{
		User:    protoUser,
		Message: "User retrieved successfully",
	}, nil
}
