package chat

import (
	"context"
	"log"
	"ws-scaling/redis"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn    *websocket.Conn
	rc      redis.Client
	channel string
}

func NewClient(conn *websocket.Conn, rc redis.Client, channel string) Client {
	return Client{
		Conn:    conn,
		rc:      rc,
		channel: channel,
	}
}

func (c Client) Reader() {
	defer c.Conn.Close()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}

		err = c.rc.Publish(context.Background(), c.channel, string(msg))
		if err != nil {
			log.Printf("Error publishing message: %v", err)
			return
		}
	}
}
