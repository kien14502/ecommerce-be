package sse

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	broker *Broker
}

func NewHandler(broker *Broker) *Handler {
	return &Handler{broker: broker}
}

func (h *Handler) Stream(c *gin.Context) {
	userID := c.GetString("user_id") // từ auth middleware

	client := &Client{
		ID:     uuid.New().String(),
		UserID: userID,
		Chan:   make(chan string, 10),
	}

	h.broker.AddClient(client)
	defer h.broker.RemoveClient(client.ID)

	// SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no") // tắt nginx buffering

	// Gửi connected event
	fmt.Fprintf(c.Writer, "event: connected\ndata: {\"status\":\"ok\"}\n\n")
	c.Writer.Flush()

	for {
		select {
		case msg, ok := <-client.Chan:
			if !ok {
				return
			}
			fmt.Fprint(c.Writer, msg)
			c.Writer.Flush()

		case <-c.Request.Context().Done():
			return // client disconnect
		}
	}
}
