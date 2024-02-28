package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/websocket"
)

type Message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

func main() {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your nickname: ")
	nickname, _ := reader.ReadString('\n')
	nickname = strings.TrimSpace(nickname)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/chat", nil)
	if err != nil {
		log.Fatal("Unable to connect to WebSocket server:", err)
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal("Unable to close connection", err)
		}
	}(conn)

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, messageBytes, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}

			var message Message
			err = json.Unmarshal(messageBytes, &message)
			if err != nil {
				log.Println("Error decoding message:", err)
				continue
			}

			fmt.Println(message.Username + ": " + message.Text)
		}
	}()

	err = conn.WriteMessage(websocket.TextMessage, []byte(nickname))
	if err != nil {
		log.Println("Error sending nickname:", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte("get_last_10"))
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}

	for {
		select {
		case <-interrupt:
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error sending close message:", err)
			}
			return
		}
	}
}