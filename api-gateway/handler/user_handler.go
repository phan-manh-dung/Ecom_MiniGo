package handler

import (
	"context"
	"net/http"
	"strconv"

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
	req := &user.GetUserRequest{Id: uint32(userID)}
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
	c.JSON(http.StatusCreated, gin.H{"user": resp.User, "message": resp.Message})
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
