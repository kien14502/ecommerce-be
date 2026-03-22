package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/internal/dto"
	"github.com/kien14502/ecommerce-be/internal/services"
	"github.com/kien14502/ecommerce-be/pkg/response"
	"github.com/kien14502/ecommerce-be/pkg/utils/cookies"
	"go.uber.org/zap"
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

	c.JSON(200, gin.H{"userID": userId})
}

// Register godoc
// @Summary      Đăng ký người dùng
// @Description  API cho phép người dùng đăng ký tài khoản bằng email và password.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "Thông tin đăng ký"
// @Success      200 {object} response.Response "Đăng ký thành công"
// @Failure      400 {object} response.Response "Dữ liệu đầu vào không hợp lệ"
// @Failure      409 {object} response.Response "Email hoặc username đã tồn tại"
// @Failure      500 {object} response.Response "Lỗi hệ thống"
// @Router       /user/register [post]
func (uc *UserController) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var in dto.RegisterRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		c.Error(response.ErrInvalidParam)
		return
	}

	err := uc.userService.Register(ctx, in)

	if err != nil {
		global.Logger.Error(
			"Register request binding failed",
			zap.Error(err),
			zap.String("path", c.FullPath()),
			zap.String("method", c.Request.Method),
			zap.String("ip", c.ClientIP()),
		)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Register successful",
	})

}

// VerifyOtp godoc
// @Summary      Xác thực OTP
// @Description  API dùng để xác thực OTP sau khi đăng ký tài khoản.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request body dto.VerifyOtpRequest true "Thông tin xác thực OTP"
// @Success      200 {object} response.Response{data=dto.LoginResponse} "Xác thực OTP thành công"
// @Failure      400 {object} response.Response "OTP không hợp lệ hoặc dữ liệu sai"
// @Failure      404 {object} response.Response "Không tìm thấy yêu cầu xác thực"
// @Failure      500 {object} response.Response "Lỗi hệ thống"
// @Router       /user/verify-otp [post]
func (uc *UserController) VerifyOtp(c *gin.Context) {
	ctx := c.Request.Context()
	var in dto.VerifyOtpRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		c.Error(response.ErrInvalidParam)
		return
	}

	res, err := uc.userService.VerifyOTP(ctx, in)
	if err != nil {
		c.Error(err)
		return
	}

	cookies.SaveRefreshToken(c.Writer, res.RefreshToken)
	c.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Verify successful",
		Data:    res,
	})
}

// VerifyOtp godoc
// @Summary      Đăng nhập tài khoản
// @Description  API dùng để xác đăng nhập tài khoản.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Thông tin đăng nhập"
// @Success      200 {object} response.Response{data=dto.LoginResponse} "Đăng nhập thành công"
// @Failure      400 {object} response.Response "Tài khoản hoặc mật khẩu không chính xác"
// @Failure      404 {object} response.Response "Không tìm thấy yêu cầu xác thực"
// @Failure      500 {object} response.Response "Lỗi hệ thống"
// @Router       /user/login [post]
func (uc *UserController) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var in dto.LoginRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		c.Error(response.ErrInvalidParam)
		return
	}
	in.DeviceID = c.GetHeader("X-Device-ID")
	if in.DeviceID == "" {
		in.DeviceID = uuid.New().String()
	}
	res, err := uc.userService.Login(ctx, in, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		c.Error(err)
	}

	cookies.SaveRefreshToken(c.Writer, res.RefreshToken)
	c.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Login successful",
		Data:    res,
	})
}

// RefreshToken godoc
// @Summary      Refresh Token
// @Description  API dùng để lấy lại refresh token và access token mới từ refresh token trong cookie
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response{data=dto.LoginResponse} "Làm mới token thành công"
// @Failure      401 {object} response.Response "Refresh token không hợp lệ hoặc đã hết hạn"
// @Failure      500 {object} response.Response "Lỗi hệ thống"
// @Router       /user/refresh-token [post]
func (uc *UserController) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	refreshToken, err := cookies.GetRefreshToken(c.Request)
	if err != nil {
		c.Error(response.ErrUnauthorized)
		return
	}
	res, err := uc.userService.RefreshToken(ctx, refreshToken)
	if err != nil {
		c.Error(response.ErrUnauthorized)
		return
	}
	c.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Success",
		Data:    res,
	})
}
