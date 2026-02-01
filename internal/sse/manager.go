// internal/sse/manager.go
package sse

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Event represents a SSE event
type Event struct {
	ID      string      `json:"id"`
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	UserID  string      `json:"user_id,omitempty"`
	Channel string      `json:"channel,omitempty"`
}

// Client represents a SSE client connection
type Client struct {
	ID         string
	UserID     string
	Channel    chan Event
	Ctx        context.Context
	CancelFunc context.CancelFunc
}

// Manager manages SSE connections and Redis pub/sub
type Manager struct {
	clients     map[string]*Client
	mu          sync.RWMutex
	redisClient *redis.Client
	ctx         context.Context
}

func NewManager(redisClient *redis.Client) *Manager {
	manager := &Manager{
		clients:     make(map[string]*Client),
		redisClient: redisClient,
		ctx:         context.Background(),
	}

	// Start listening to Redis pub/sub
	go manager.listenToRedis()

	return manager
}

// AddClient adds a new SSE client
func (m *Manager) AddClient(clientID, userID string) *Client {
	ctx, cancel := context.WithCancel(m.ctx)

	client := &Client{
		ID:         clientID,
		UserID:     userID,
		Channel:    make(chan Event, 100), // Buffer để tránh block
		Ctx:        ctx,
		CancelFunc: cancel,
	}

	m.mu.Lock()
	m.clients[clientID] = client
	m.mu.Unlock()

	log.Printf("Client connected: %s (User: %s)", clientID, userID)
	return client
}

// RemoveClient removes a SSE client
func (m *Manager) RemoveClient(clientID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if client, exists := m.clients[clientID]; exists {
		client.CancelFunc()
		close(client.Channel)
		delete(m.clients, clientID)
		log.Printf("Client disconnected: %s", clientID)
	}
}

// PublishToRedis publishes an event to Redis
func (m *Manager) PublishToRedis(channel string, event Event) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = m.redisClient.Publish(m.ctx, channel, eventJSON).Err()
	if err != nil {
		return fmt.Errorf("failed to publish to redis: %w", err)
	}

	log.Printf("Published to Redis channel '%s': %s", channel, event.Type)
	return nil
}

// listenToRedis subscribes to Redis channels and broadcasts to clients
func (m *Manager) listenToRedis() {
	// Subscribe to multiple channels
	pubsub := m.redisClient.PSubscribe(m.ctx,
		"notifications:*", // All notification channels
		"user:*",          // User-specific channels
		"broadcast",       // Global broadcast channel
	)
	defer pubsub.Close()

	ch := pubsub.Channel()

	log.Println("Started listening to Redis pub/sub")

	for msg := range ch {
		var event Event
		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			log.Printf("Error unmarshaling event: %v", err)
			continue
		}

		event.Channel = msg.Channel
		m.broadcastEvent(event)
	}
}

// broadcastEvent sends event to relevant clients
func (m *Manager) broadcastEvent(event Event) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, client := range m.clients {
		// Broadcast to all if it's a global event
		if event.Channel == "broadcast" {
			m.sendToClient(client, event)
			continue
		}

		// Send to specific user
		if event.UserID != "" && event.UserID == client.UserID {
			m.sendToClient(client, event)
			continue
		}

		// Send to clients listening to specific channel
		if event.Channel != "" {
			m.sendToClient(client, event)
		}
	}
}

// sendToClient sends event to a specific client
func (m *Manager) sendToClient(client *Client, event Event) {
	select {
	case client.Channel <- event:
		log.Printf("Sent event '%s' to client %s", event.Type, client.ID)
	case <-time.After(1 * time.Second):
		log.Printf("Timeout sending to client %s", client.ID)
	case <-client.Ctx.Done():
		log.Printf("Client %s context done", client.ID)
	}
}

// GetClientCount returns the number of connected clients
func (m *Manager) GetClientCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.clients)
}
