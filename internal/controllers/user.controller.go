package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kien14502/ecommerce-be/internal/dto"
	"github.com/kien14502/ecommerce-be/internal/services"
	"github.com/kien14502/ecommerce-be/pkg/response"
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

// Register godoc
// @Summary      Đăng ký người dùng mới
// @Description  API cho phép người dùng đăng ký tài khoản bằng Email và Password.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request body      models.RegisterRequest  true  "Thông tin đăng ký (Email, Password)"
// @Success      200     {object}  map[string]string       "Trả về message thành công"
// @Failure      400     {object}  map[string]string       "Lỗi dữ liệu đầu vào không hợp lệ"
// @Router       /user/register [post]
func (uc *UserController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	status := uc.userService.Register(req.Email, req.Password)

	c.JSON(http.StatusOK, gin.H{
		"message": response.GetMessage(status),
	})
}
