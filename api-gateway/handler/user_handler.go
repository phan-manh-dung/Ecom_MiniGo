package handler

import (
	"context"
	"net/http"
	"strconv"

	"gin/api-gateway/service_manager"
	"gin/api-gateway/utils"
	"gin/proto/generated/user"

	"github.com/gin-gonic/gin"
)

type UserServiceClient struct {
	ServiceManager *service_manager.ServiceManager
}

func NewUserServiceClient(serviceManager *service_manager.ServiceManager) *UserServiceClient {
	return &UserServiceClient{ServiceManager: serviceManager}
}

func (u *UserServiceClient) GetUser(c *gin.Context) {
	if u.ServiceManager == nil || u.ServiceManager.UserClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "User service not available"})
		return
	}

	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	// tạo request
	req := &user.GetUserRequest{Id: uint32(userID)}
	// gọi grpc
	resp, err := u.ServiceManager.UserClient.GetUser(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": resp.User, "message": resp.Message})
}

func (u *UserServiceClient) GetUserBySDT(c *gin.Context) {
	if u.ServiceManager == nil || u.ServiceManager.UserClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "User service not available"})
		return
	}

	sdt := c.Param("id")
	req := &user.GetUserBySDTRequest{Sdt: sdt}
	resp, err := u.ServiceManager.UserClient.GetUserBySDT(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": resp.User, "message": resp.Message})
}

func (u *UserServiceClient) CreateUser(c *gin.Context) {
	if u.ServiceManager == nil || u.ServiceManager.UserClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "User service not available"})
		return
	}

	var createReq struct {
		Name   string `json:"name" binding:"required"`
		SDT    string `json:"sdt" binding:"required"`
		RoleID uint32 `json:"role_id"`
	}
	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	req := &user.CreateUserRequest{
		Name:   createReq.Name,
		Sdt:    createReq.SDT,
		RoleId: createReq.RoleID,
	}
	resp, err := u.ServiceManager.UserClient.CreateUser(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Kiểm tra user có account và role không
	if resp.User == nil || resp.User.Account == nil || resp.User.Account.Role == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User account or role not found"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":    resp.User,
		"message": resp.Message,
	})
}

func (u *UserServiceClient) UpdateUser(c *gin.Context) {
	if u.ServiceManager == nil || u.ServiceManager.UserClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "User service not available"})
		return
	}

	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var updateReq struct {
		Name string `json:"name"`
		SDT  string `json:"sdt"`
	}
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	req := &user.UpdateUserRequest{
		Id:   uint32(userID),
		Name: updateReq.Name,
		Sdt:  updateReq.SDT,
	}
	resp, err := u.ServiceManager.UserClient.UpdateUser(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": resp.User, "message": resp.Message})
}

func (u *UserServiceClient) DeleteUser(c *gin.Context) {
	if u.ServiceManager == nil || u.ServiceManager.UserClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "User service not available"})
		return
	}

	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	req := &user.DeleteUserRequest{Id: uint32(userID)}
	resp, err := u.ServiceManager.UserClient.DeleteUser(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": resp.Message})
}

func (u *UserServiceClient) ListUsers(c *gin.Context) {
	if u.ServiceManager == nil || u.ServiceManager.UserClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "User service not available"})
		return
	}

	req := &user.ListUsersRequest{}
	resp, err := u.ServiceManager.UserClient.ListUsers(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": resp.Users, "message": resp.Message})
}

func (u *UserServiceClient) Login(c *gin.Context) {
	if u.ServiceManager == nil || u.ServiceManager.UserClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "User service not available"})
		return
	}

	var loginReq struct {
		SDT string `json:"sdt" binding:"required"`
	}
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	req := &user.GetUserBySDTRequest{Sdt: loginReq.SDT}
	resp, err := u.ServiceManager.UserClient.GetUserBySDT(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Kiểm tra user có tồn tại không
	if resp.User == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Kiểm tra user có account và role không
	if resp.User.Account == nil || resp.User.Account.Role == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User account or role not found"})
		return
	}

	// Tạo JWT token
	roleName := resp.User.Account.Role.Name
	token, err := utils.GenerateToken(strconv.FormatUint(uint64(resp.User.Id), 10), resp.User.Name, roleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    resp.User,
		"message": "Login successful",
		"token":   token,
	})
}
