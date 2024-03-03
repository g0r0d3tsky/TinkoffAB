package main

import (
	"2024-spring-ab-go-hw-1-template-g0r0d3tsky/config"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/websocket"
)

type Message struct {
	UserNickname string
	Content      string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	c, err := config.Read()

	if err != nil {
		log.Println("failed to read config:", err.Error())
		return
	}
	conn, _, err := websocket.DefaultDialer.Dial(c.ServerAddress("chat"), nil) //вынести в config
	if err != nil {
		log.Fatal("Unable to connect to WebSocket server:", err)
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal("Unable to close connection", err)
		}
	}(conn)

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your nickname: ")
	nickname, _ := reader.ReadString('\n')
	nickname = strings.TrimSpace(nickname)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go readMessages(ctx, conn, done)

	go writeMessages(conn, reader, nickname)

	select {
	case <-interrupt:
		fmt.Println("Interrupt signal received. Shutting down...")
		cancel()
	}

	<-done
	fmt.Println("Exiting...")
}

func readMessages(ctx context.Context, conn *websocket.Conn, done chan struct{}) {
	defer close(done)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, messageBytes, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}

			var receivedMessage Message
			err = json.Unmarshal(messageBytes, &receivedMessage)
			if err != nil {
				log.Println("Error decoding message:", err)
				continue
			}

			fmt.Println(receivedMessage.UserNickname + ": " + receivedMessage.Content)
		}
	}
}

func writeMessages(conn *websocket.Conn, reader *bufio.Reader, nickname string) {
	for {
		content, _ := reader.ReadString('\n')
		content = strings.TrimSpace(content)

		message := Message{
			UserNickname: nickname,
			Content:      content,
		}

		messageBytes, err := json.Marshal(message)
		if err != nil {
			log.Println("Error encoding message:", err)
			continue
		}

		err = conn.WriteMessage(websocket.TextMessage, messageBytes)
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
	}
}
