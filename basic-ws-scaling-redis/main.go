package main

import (
	"context"
	"log"
	"net/http"
	"ws-scaling/redis"

	"github.com/gorilla/websocket"
)

const channel = "chat"

var ugrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type client struct {
	conn *websocket.Conn
}

func (c client) reader() {
	defer c.conn.Close()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}

		err = redisClient.Publish(context.Background(), channel, string(msg))
		if err != nil {
			log.Printf("Error publishing message: %v", err)
			return
		}
	}
}

var clients []client
var redisClient redis.Client

func chatHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := ugrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Error upgrading to websocket: %v", err)
		return
	}

	client := client{
		conn: conn,
	}

	clients = append(clients, client)

	go client.reader()

	redisClient.Publish(context.Background(), channel, "New user joined")
}

func main() {
	clients = make([]client, 0)

	redisClient = redis.NewClient()

	subscriber := redisClient.Subscribe(context.Background(), channel)

	go func() {
		for {
			msg, err := subscriber.ReceiveMessage(context.Background())
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				return
			}

			for _, client := range clients {
				err = client.conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
				if err != nil {
					log.Printf("Error sending message: %v", err)
					return
				}
			}
		}
	}()

	http.HandleFunc("/chat", chatHandler)

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
