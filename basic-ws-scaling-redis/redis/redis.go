package redis

import (
	"context"
	"fmt"
	"os"

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
	var host, port string

	if os.Getenv("REDIS_HOST") != "" {
		host = os.Getenv("REDIS_HOST")
	} else {
		host = "localhost"
	}

	if os.Getenv("REDIS_PORT") != "" {
		port = os.Getenv("REDIS_PORT")
	} else {
		port = "6379"
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	return &clientImpl{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (c *clientImpl) Publish(ctx context.Context, channel string, message interface{}) error {
	return c.client.Publish(ctx, channel, message).Err()
}

func (c *clientImpl) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return c.client.Subscribe(ctx, channel)
}
