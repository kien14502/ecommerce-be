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

// GetUser godoc
// @Summary      Get user information
// @Description  Retrieve user details by ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  object{userID=string}  "User information"
// @Failure      400  {object}  object{error=string}   "Bad request"
// @Failure      404  {object}  object{error=string}   "User not found"
// @Router       /users/{id} [get]
func (uc *UserController) GetUser(c *gin.Context) {
	userId := c.Param("id")
	// userId := uc.userService.GetUserName("123")

	c.JSON(200, gin.H{"userID": userId})
}
