// internal/transport/sse/broker.go
package sse

import (
	"context"
	"fmt"
	"sync"

	"github.com/kien14502/ecommerce-be/internal/transport/pubsub"
)

type Client struct {
	ID     string
	UserID string
	Chan   chan string
}

type Broker struct {
	clients map[string]*Client
	mu      sync.RWMutex
	pubsub  *pubsub.RedisPubSub
}

func NewBroker(ps *pubsub.RedisPubSub) *Broker {
	b := &Broker{
		clients: make(map[string]*Client),
		pubsub:  ps,
	}
	// Subscribe inventory updates cho tất cả clients
	go b.listenInventory()
	return b
}

func (b *Broker) AddClient(client *Client) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clients[client.ID] = client

	// Subscribe notification riêng cho user này
	go b.listenNotification(client)
}

func (b *Broker) RemoveClient(clientID string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if client, ok := b.clients[clientID]; ok {
		close(client.Chan)
		delete(b.clients, clientID)
	}
}

func (b *Broker) listenInventory() {
	sub := b.pubsub.Subscribe(context.Background(), pubsub.ChannelInventory)
	defer sub.Close()

	for msg := range sub.Channel() {
		// Broadcast tới tất cả clients
		b.mu.RLock()
		for _, client := range b.clients {
			select {
			case client.Chan <- fmt.Sprintf("event: inventory\ndata: %s\n\n", msg.Payload):
			default: // client slow/disconnected, skip
			}
		}
		b.mu.RUnlock()
	}
}

func (b *Broker) listenNotification(client *Client) {
	channel := fmt.Sprintf(pubsub.ChannelNotification, client.UserID)
	sub := b.pubsub.Subscribe(context.Background(), channel)
	defer sub.Close()

	for msg := range sub.Channel() {
		select {
		case client.Chan <- fmt.Sprintf("event: notification\ndata: %s\n\n", msg.Payload):
		default:
			return // client disconnected
		}
	}
}
