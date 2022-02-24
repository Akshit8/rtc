package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/"}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial error:", err)
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			log.Println("Recieved:", string(msg))
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Hello from client at %v", t)))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interupt:
			log.Println("interupt")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
