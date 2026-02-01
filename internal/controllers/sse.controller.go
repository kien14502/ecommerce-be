// internal/controller/sse_controller.go
package controllers

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kien14502/ecommerce-be/internal/sse"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SSEController struct {
	manager *sse.Manager
}

func NewSSEController(manager *sse.Manager) *SSEController {
	return &SSEController{
		manager: manager,
	}
}

// StreamEvents godoc
// @Summary      Stream Server-Sent Events
// @Description  Establishes an SSE connection to receive real-time events for a specific user
// @Tags         SSE
// @Accept       json
// @Produce      text/event-stream
// @Param        user_id  query  string  false  "User ID (defaults to 'anonymous' if not provided)"
// @Success      200  {object}  object  "SSE connection established. Events will be streamed in the format: id, event, data"
// @Header       200  {string}  Content-Type         "text/event-stream"
// @Header       200  {string}  Cache-Control        "no-cache"
// @Header       200  {string}  Connection           "keep-alive"
// @Header       200  {string}  Transfer-Encoding    "chunked"
// @Header       200  {string}  X-Accel-Buffering    "no"
// @Router       /sse/events [get]
func (ctrl *SSEController) StreamEvents(c *gin.Context) {
	// Get user ID from auth (hoặc từ query param để test)
	userID := c.Query("user_id")
	if userID == "" {
		userID = "anonymous"
	}

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")
	c.Header("X-Accel-Buffering", "no") // Disable nginx buffering

	// Create client
	clientID := uuid.New().String()
	client := ctrl.manager.AddClient(clientID, userID)
	defer ctrl.manager.RemoveClient(clientID)

	// Send initial connection event
	ctrl.sendSSE(c.Writer, sse.Event{
		ID:   uuid.New().String(),
		Type: "connected",
		Data: map[string]interface{}{
			"client_id": clientID,
			"message":   "Connected to event stream",
		},
	})

	// Keep-alive ticker
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Stream events
	for {
		select {
		case <-c.Request.Context().Done():
			return

		case event := <-client.Channel:
			if err := ctrl.sendSSE(c.Writer, event); err != nil {
				return
			}

		case <-ticker.C:
			// Send heartbeat
			ctrl.sendSSE(c.Writer, sse.Event{
				Type: "heartbeat",
				Data: map[string]interface{}{"timestamp": time.Now().Unix()},
			})
		}
	}
}

// sendSSE sends a SSE formatted message
func (ctrl *SSEController) sendSSE(w io.Writer, event sse.Event) error {
	// Format: id, event, data
	if event.ID != "" {
		fmt.Fprintf(w, "id: %s\n", event.ID)
	}

	if event.Type != "" {
		fmt.Fprintf(w, "event: %s\n", event.Type)
	}

	fmt.Fprintf(w, "data: %v\n\n", formatData(event.Data))

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	return nil
}

func formatData(data interface{}) string {
	switch v := data.(type) {
	case string:
		return v
	default:
		// Convert to JSON string
		return fmt.Sprintf("%v", v)
	}
}

// PublishEvent publishes an event via Redis
func (ctrl *SSEController) PublishEvent(c *gin.Context) {
	var req struct {
		Channel string      `json:"channel" binding:"required"`
		Type    string      `json:"type" binding:"required"`
		Data    interface{} `json:"data" binding:"required"`
		UserID  string      `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := sse.Event{
		ID:      uuid.New().String(),
		Type:    req.Type,
		Data:    req.Data,
		UserID:  req.UserID,
		Channel: req.Channel,
	}

	if err := ctrl.manager.PublishToRedis(req.Channel, event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Event published",
		"event":   event,
	})
}

// GetStats returns SSE statistics
func (ctrl *SSEController) GetStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"connected_clients": ctrl.manager.GetClientCount(),
		"timestamp":         time.Now().Unix(),
	})
}
