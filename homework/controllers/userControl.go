package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mcyuchen915/golang_gin_gorm_homework.git/homework/homework/models"
	"github.com/mcyuchen915/golang_gin_gorm_homework.git/homework/homework/services"
	"github.com/mcyuchen915/golang_gin_gorm_homework.git/homework/homework/utils"
)
//用户控制器
type UserController struct{}

var UserService *services.UserService

// Register 用户注册
func (uc *UserController) Register(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	user, err := UserService.RegisterUser(req)
	if err != nil {
		utils.InternalServerError(c, "Failed to create user")
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.InternalServerError(c, "Failed to generate token")
		return
	}

	utils.Success(c, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// Login 用户登录
func (ac *UserController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 查找用户
	user, err := UserService.Authenticate(req.Username, req.Password)
	if err != nil {
		utils.Unauthorized(c, "Invalid username or password")
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.InternalServerError(c, "Failed to generate token")
		return
	}

	utils.Success(c, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// GetUser 获取用户信息
func (ac *UserController) GetUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}
	user, err := UserService.GetUserByID(userID.(uint))
	if err != nil {
		utils.Unauthorized(c, "User not authenticated")
		return
	}
	utils.Success(c, user)
}
