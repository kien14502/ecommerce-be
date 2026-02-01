package services

import (
	"github.com/kien14502/ecommerce-be/internal/sse"

	"github.com/google/uuid"
)

type NotificationService struct {
	sseManager *sse.Manager
}

func NewNotificationService(sseManager *sse.Manager) *NotificationService {
	return &NotificationService{
		sseManager: sseManager,
	}
}

// NotifyUser sends notification to specific user
func (s *NotificationService) NotifyUser(userID string, notifType string, data interface{}) error {
	event := sse.Event{
		ID:     uuid.New().String(),
		Type:   notifType,
		Data:   data,
		UserID: userID,
	}

	return s.sseManager.PublishToRedis("user:"+userID, event)
}

// BroadcastNotification sends notification to all connected clients
func (s *NotificationService) BroadcastNotification(notifType string, data interface{}) error {
	event := sse.Event{
		ID:   uuid.New().String(),
		Type: notifType,
		Data: data,
	}

	return s.sseManager.PublishToRedis("broadcast", event)
}

// NotifyChannel sends notification to a specific channel
func (s *NotificationService) NotifyChannel(channel string, notifType string, data interface{}) error {
	event := sse.Event{
		ID:   uuid.New().String(),
		Type: notifType,
		Data: data,
	}

	return s.sseManager.PublishToRedis("notifications:"+channel, event)
}
