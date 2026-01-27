package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/kien14502/ecommerce-be/internal/services"
)

type UserController struct {
	userService services.IUserService
}

func NewUserController(userService services.IUserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (uc *UserController) GetUser(c *gin.Context) {
	userId := uc.userService.GetUserName("123")

	c.JSON(200, gin.H{"userID": userId})
}
