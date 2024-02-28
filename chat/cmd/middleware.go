package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func (app *app) SendLastMessagesMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		connection, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading to WebSocket connection:", err)
			return
		}
		defer connection.Close()

		ctx := r.Context()
		app.sendLastMessages(ctx, connection)

		next.ServeHTTP(w, r)
	})
}

func (app *app) sendLastMessages(ctx context.Context, connection *websocket.Conn) {
	messages, err := app.UC.GetAmountMessage(ctx, 10)
	if err != nil {
		log.Fatal("getting messages: ", err)
		return
	}
	for _, msg := range messages {
		messageBytes, err := json.Marshal(msg)
		if err != nil {
			log.Println("Error encoding message:", err)
			continue
		}

		err = connection.WriteMessage(websocket.TextMessage, messageBytes)
		if err != nil {
			log.Println("Error sending message:", err)
			continue
		}
	}
}