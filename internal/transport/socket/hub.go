package socket

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kien14502/ecommerce-be/internal/transport/pubsub"
)

type Hub struct {
	rooms      map[string]map[*Client]bool // roomID → clients
	register   chan *Client
	unregister chan *Client
	publish    chan ChatMessage
	pubsub     *pubsub.RedisPubSub
}

func NewHub(ps *pubsub.RedisPubSub) *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		publish:    make(chan ChatMessage),
		pubsub:     ps,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if h.rooms[client.RoomID] == nil {
				h.rooms[client.RoomID] = make(map[*Client]bool)
				go h.listenRoom(client.RoomID) // subscribe Redis channel cho room này
			}
			h.rooms[client.RoomID][client] = true

		case client := <-h.unregister:
			if room, ok := h.rooms[client.RoomID]; ok {
				delete(room, client)
				close(client.send)
			}

		case msg := <-h.publish:
			// Publish lên Redis → sync tới tất cả server instances
			h.pubsub.Publish(context.Background(),
				fmt.Sprintf(pubsub.ChannelChat, msg.RoomID),
				pubsub.Message{
					Type:    "chat",
					Payload: mustMarshal(msg),
					RoomID:  msg.RoomID,
				},
			)
		}
	}
}

func (h *Hub) listenRoom(roomID string) {
	channel := fmt.Sprintf(pubsub.ChannelChat, roomID)
	sub := h.pubsub.Subscribe(context.Background(), channel)
	defer sub.Close()

	for redisMsg := range sub.Channel() {
		var msg pubsub.Message
		json.Unmarshal([]byte(redisMsg.Payload), &msg)

		// Broadcast tới tất cả clients trong room này
		if room, ok := h.rooms[roomID]; ok {
			for client := range room {
				select {
				case client.send <- msg.Payload:
				default:
					close(client.send)
					delete(room, client)
				}
			}
		}
	}
}

func mustMarshal(v any) json.RawMessage {
	data, _ := json.Marshal(v)
	return data
}
