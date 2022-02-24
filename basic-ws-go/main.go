package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleWebSocket(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	signal := make(chan bool)

	conn.SetCloseHandler(func(code int, text string) error {
		log.Println("Got close message:", code, text)
		signal <- true
		conn.Close()
		return nil
	})

	go func() {
		for {
			select {
			case <-signal:
				return
			default:
				mt, msg, err := conn.ReadMessage()
				if err != nil {
					log.Println("Error reading message:", err)
					return
				}
				log.Printf("Got message: %s", msg)

				err = conn.WriteMessage(mt, msg)
				if err != nil {
					log.Println("Error writing message:", err)
					return
				}
			}
		}
	}()
}

func WriteMessageEvery5s(conn *websocket.Conn) {
	for {
		err := conn.WriteMessage(websocket.TextMessage, []byte("Hello from server"))
		if err != nil {
			log.Println("Error writing message:", err)
			return
		}

		time.Sleep(5 * time.Second)
	}
}

func main() {
	http.HandleFunc("/", handleWebSocket)

	log.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", nil)
}
