package middlewares

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kien14502/ecommerce-be/pkg/response"
)

type AppError struct {
	Status  int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var appErr *response.AppError
		if errors.As(err, &appErr) {
			c.JSON(appErr.HTTPStatus, response.Response{
				Success: false,
				Code:    appErr.Code,
				Message: appErr.Message,
			})
			c.Abort()
			return
		}

		// fallback
		fmt.Println("Unhandled error:", err)

		c.JSON(500, response.Response{
			Success: false,
			Code:    "SYS0001",
			Message: "Internal server error",
		})
		c.Abort()
	}
}
