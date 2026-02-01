package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		err := c.Errors[0].Err

		if appErr, ok := err.(*AppError); ok {
			c.AbortWithStatusJSON(appErr.Status, gin.H{
				"success": false,
				"message": appErr.Message,
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
	}
}
