package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"ws-scaling/chat"
	"ws-scaling/redis"

	"github.com/gorilla/websocket"
)

const channel = "chat"

var ugrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients []chat.Client
var redisClient redis.Client
var apptag string

func chatHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := ugrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Error upgrading to websocket: %v", err)
		return
	}

	conn.SetCloseHandler(func(code int, text string) error {
		log.Println("removing client")
		for i, client := range clients {
			if client.Conn == conn {
				clients = append(clients[:i], clients[i+1:]...)
				break
			}
		}

		return nil
	})

	client := chat.NewClient(conn, redisClient, channel)

	clients = append(clients, client)

	go client.Reader()

	redisClient.Publish(context.Background(), channel, "new user joined the chat")
}

func main() {
	if os.Getenv("APP_TAG") != "" {
		apptag = os.Getenv("APP_TAG")
	} else {
		apptag = "chat"
	}

	clients = make([]chat.Client, 0)

	redisClient = redis.NewClient()

	subscriber := redisClient.Subscribe(context.Background(), channel)

	ch := subscriber.Channel()

	go func() {
		for msg := range ch {
			// msg, err := subscriber.ReceiveMessage(context.Background())
			// if err != nil {
			// 	log.Printf("Error receiving message: %v", err)
			// 	return
			// }

			for _, client := range clients {
				newMsg := fmt.Sprintf("%s: %s", apptag, msg.Payload)
				err := client.Conn.WriteMessage(websocket.TextMessage, []byte(newMsg))
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
