package main

import (
	"2024-spring-ab-go-hw-1-template-g0r0d3tsky/chat/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"log/slog"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {

			return true
		},
	}

	clients = make(map[*websocket.Conn]struct{})
)

func (app *app) echo(w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)
	defer func(connection *websocket.Conn) {
		err := connection.Close()
		if err != nil {
			slog.Error("connection closing", err)
		}
	}(connection)

	clients[connection] = struct{}{}
	defer delete(clients, connection)

	for {
		mt, messageBytes, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break
		}

		var message *domain.Message
		err = json.Unmarshal(messageBytes, &message)
		if err != nil {
			log.Println("Error decoding message:", err)
			continue
		}
		err = app.UC.CreateMessage(context.TODO(), message)
		if err != nil {
			log.Println("Error saving message:", err)
			return
		}

		// Теперь мы рассылаем сообщения всем клиентам
		go writeMessage(message)

		go messageHandler(message)
	}
}

func writeMessage(message *domain.Message) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Println("Error encoding message:", err)
		return
	}

	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, messageBytes)
		if err != nil {
			return
		}
	}
}

func messageHandler(message *domain.Message) {
	fmt.Printf("%v : , %v", message.UserNickname, message.Content)

}
func (app *app) RegisterHandlers() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/chat", app.SendLastMessagesMiddleware(app.echo)).Methods("GET")
	return router
}