package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// HTTP Status Messages
var httpStatusMessages = map[int]string{
	http.StatusBadRequest:          "Bad Request",
	http.StatusUnauthorized:        "Unauthorized",
	http.StatusForbidden:           "Forbidden",
	http.StatusNotFound:            "Not Found",
	http.StatusConflict:            "Conflict",
	http.StatusTooManyRequests:     "Too Many Requests",
	http.StatusInternalServerError: "Internal Server Error",
	http.StatusServiceUnavailable:  "Service Unavailable",
	http.StatusGatewayTimeout:      "Gateway Timeout",
}

// GetHTTPStatusMessage returns message for HTTP status code
func GetHTTPStatusMessage(statusCode int) string {
	if msg, ok := httpStatusMessages[statusCode]; ok {
		return msg
	}
	return "Unknown Error"
}

// SuccessResponse sends success response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends error response
func ErrorResponse(c *gin.Context, statusCode int, reason string) {
	errorMsg := reason
	if errorMsg == "" {
		errorMsg = GetHTTPStatusMessage(statusCode)
	}

	c.JSON(statusCode, Response{
		Success: false,
		Error:   errorMsg,
	})
}

// ============ Helper Functions ============

// Success (200 OK)
func Success(c *gin.Context, data interface{}) {
	SuccessResponse(c, http.StatusOK, "Success", data)
}

// SuccessWithMessage (200 OK with custom message)
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, http.StatusOK, message, data)
}

// Created (201)
func Created(c *gin.Context, data interface{}) {
	SuccessResponse(c, http.StatusCreated, "Created successfully", data)
}

// NoContent (204)
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// BadRequest (400)
func BadRequest(c *gin.Context, reason string) {
	ErrorResponse(c, http.StatusBadRequest, reason)
}

// Unauthorized (401)
func Unauthorized(c *gin.Context, reason string) {
	ErrorResponse(c, http.StatusUnauthorized, reason)
}

// Forbidden (403)
func Forbidden(c *gin.Context, reason string) {
	ErrorResponse(c, http.StatusForbidden, reason)
}

// NotFound (404)
func NotFound(c *gin.Context, reason string) {
	ErrorResponse(c, http.StatusNotFound, reason)
}

// Conflict (409)
func Conflict(c *gin.Context, reason string) {
	ErrorResponse(c, http.StatusConflict, reason)
}

// TooManyRequests (429)
func TooManyRequests(c *gin.Context, reason string) {
	ErrorResponse(c, http.StatusTooManyRequests, reason)
}

// InternalServerError (500)
func InternalServerError(c *gin.Context, reason string) {
	ErrorResponse(c, http.StatusInternalServerError, reason)
}

// ServiceUnavailable (503)
func ServiceUnavailable(c *gin.Context, reason string) {
	ErrorResponse(c, http.StatusServiceUnavailable, reason)
}
