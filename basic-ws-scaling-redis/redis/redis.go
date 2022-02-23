package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Client interface {
	Publish(ctx context.Context, channel string, message interface{}) error
	Subscribe(ctx context.Context, channel string) *redis.PubSub
}

type clientImpl struct {
	client *redis.Client
}

func NewClient() Client {
	return &clientImpl{
		client: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	}
}

func (c *clientImpl) Publish(ctx context.Context, channel string, message interface{}) error {
	return c.client.Publish(ctx, channel, message).Err()
}

func (c *clientImpl) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return c.client.Subscribe(ctx, channel)
}
