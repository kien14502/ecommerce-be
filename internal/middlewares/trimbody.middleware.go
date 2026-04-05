package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kien14502/ecommerce-be/pkg/utils"
)

func TrimBodyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if !strings.Contains(c.GetHeader("Content-Type"), "application/json") {
			c.Next()
			return
		}

		if c.Request.Body == nil {
			c.Next()
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Failed to read request body",
			})
			return
		}
		defer c.Request.Body.Close()

		var data interface{}
		if err := json.Unmarshal(body, &data); err != nil {

			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			c.Next()
			return
		}

		trimmed := utils.TrimValue(data)

		trimmedBody, err := json.Marshal(trimmed)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to process request body",
			})
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(trimmedBody))
		c.Request.ContentLength = int64(len(trimmedBody))

		c.Next()
	}
}
