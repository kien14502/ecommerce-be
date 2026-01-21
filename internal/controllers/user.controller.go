package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/kien14502/ecommerce-be/internal/services"
)

type UserController struct{
	userService *services.UserServiceType
}

func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

func (uc *UserController) GetUser(c *gin.Context) {
	c.JSON(200, gin.H{
		"user": uc.userService.GetUserName("123"),
	})
	
}