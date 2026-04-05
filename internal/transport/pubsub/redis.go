package pubsub

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	UserID  string          `json:"user_id,omitempty"`
	RoomID  string          `json:"room_id,omitempty"`
}

type RedisPubSub struct {
	client *redis.Client
}

func NewRedisPubSub(client *redis.Client) *RedisPubSub {
	return &RedisPubSub{client: client}
}

func (r *RedisPubSub) Publish(ctx context.Context, channel string, msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return r.client.Publish(ctx, channel, data).Err()
}

func (r *RedisPubSub) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return r.client.Subscribe(ctx, channels...)
}

// Channel naming convention
const (
	ChannelInventory    = "inventory:updates"
	ChannelNotification = "notification:user:%s" // notification:user:<userID>
	ChannelChat         = "chat:room:%s"         // chat:room:<roomID>
)
