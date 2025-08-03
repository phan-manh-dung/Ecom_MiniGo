package handler

import (
	"context"
	"net/http"
	"strconv"

	"gin/api-gateway/utils"
	"gin/proto/generated/user"

	"github.com/gin-gonic/gin"
)

type UserServiceClient struct {
	Client user.UserServiceClient
}

func NewUserServiceClient(client user.UserServiceClient) *UserServiceClient {
	return &UserServiceClient{Client: client}
}

func (u *UserServiceClient) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	// tạo request
	req := &user.GetUserRequest{Id: uint32(userID)}
	// gọi grpc
	resp, err := u.Client.GetUser(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": resp.User, "message": resp.Message})
}

func (u *UserServiceClient) GetUserBySDT(c *gin.Context) {
	sdt := c.Param("sdt")
	req := &user.GetUserBySDTRequest{Sdt: sdt}
	resp, err := u.Client.GetUserBySDT(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": resp.User, "message": resp.Message})
}

func (u *UserServiceClient) CreateUser(c *gin.Context) {
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
	resp, err := u.Client.CreateUser(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Tạo JWT token cho user mới - lấy role từ database
	roleName := resp.User.Account.Role.Name
	token, err := utils.GenerateToken(strconv.FormatUint(uint64(resp.User.Id), 10), resp.User.Name, roleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":    resp.User,
		"message": resp.Message,
		"token":   token,
	})
}

func (u *UserServiceClient) UpdateUser(c *gin.Context) {
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
	resp, err := u.Client.UpdateUser(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": resp.User, "message": resp.Message})
}

func (u *UserServiceClient) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	req := &user.DeleteUserRequest{Id: uint32(userID)}
	resp, err := u.Client.DeleteUser(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": resp.Message})
}

func (u *UserServiceClient) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	req := &user.ListUsersRequest{Page: int32(page), Limit: int32(limit)}
	resp, err := u.Client.ListUsers(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": resp.Users, "total": resp.Total, "message": resp.Message})
}

// LoginUser xử lý đăng nhập và trả về JWT token
func (u *UserServiceClient) LoginUser(c *gin.Context) {
	var loginReq struct {
		SDT string `json:"sdt" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Lấy thông tin user theo SDT
	req := &user.GetUserBySDTRequest{Sdt: loginReq.SDT}
	resp, err := u.Client.GetUserBySDT(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Kiểm tra user có tồn tại không
	if resp.User == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Tạo JWT token - lấy role từ database
	var roleName string
	if resp.User.Account != nil && resp.User.Account.Role != nil {
		roleName = resp.User.Account.Role.Name
	} else {
		roleName = "USER" // default role
	}
	token, err := utils.GenerateToken(strconv.FormatUint(uint64(resp.User.Id), 10), resp.User.Name, roleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    resp.User,
		"token":   token,
		"message": "Login successful",
	})
}
