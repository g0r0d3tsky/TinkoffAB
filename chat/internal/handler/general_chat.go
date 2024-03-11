package handler

import (
	"2024-spring-ab-go-hw-1-template-g0r0d3tsky/chat/internal/domain"
	"2024-spring-ab-go-hw-1-template-g0r0d3tsky/chat/internal/handler/models"
	"2024-spring-ab-go-hw-1-template-g0r0d3tsky/chat/internal/usecase"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"log/slog"
	"net/http"
	"sync"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {

			return true
		},
	}
)

func (h *Handler) AddClient(client *websocket.Conn) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.Clients[client] = struct{}{}
}

func (h *Handler) RemoveClient(client *websocket.Conn) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.Clients, client)
}

func (h *Handler) GetClients() map[*websocket.Conn]struct{} {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	clientsCopy := make(map[*websocket.Conn]struct{})
	for client := range h.Clients {
		clientsCopy[client] = struct{}{}
	}
	return clientsCopy
}

func NewHandler(services usecase.MessageUseCase) *Handler {
	return &Handler{
		Services: services,
		Clients:  make(map[*websocket.Conn]struct{}),
		mutex:    sync.Mutex{},
	}
}

func (h *Handler) Echo(w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)
	defer func(connection *websocket.Conn) {
		err := connection.Close()
		if err != nil {
			slog.Error("connection closing", err)
		}
	}(connection)

	h.mutex.Lock()
	h.Clients[connection] = struct{}{}
	h.mutex.Unlock()

	defer func() {
		h.mutex.Lock()
		delete(h.Clients, connection)
		h.mutex.Unlock()
	}()

	h.sendLastMessages(r.Context(), connection)

	for {
		mt, messageBytes, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break
		}

		var message models.Message
		err = json.Unmarshal(messageBytes, &message)
		if err != nil {
			log.Println("Error decoding message:", err)
			continue
		}
		m := &domain.Message{
			UserNickname: message.UserNickname,
			Content:      message.Content,
		}
		err = h.Services.CreateMessage(r.Context(), m)
		if err != nil {
			log.Println("Error saving message:", err)
			return
		}

		// Теперь мы рассылаем сообщения всем клиентам
		go h.writeMessage(r.Context(), message)
		go h.MessageHandler(message)
	}
}

func (h *Handler) writeMessage(ctx context.Context, message models.Message) {
	select {
	case <-ctx.Done():
		return
	default:
		messageBytes, err := json.Marshal(message)
		if err != nil {
			log.Println("Error encoding message:", err)
			return
		}

		h.mutex.Lock()
		defer h.mutex.Unlock()

		for conn := range h.Clients {
			err := conn.WriteMessage(websocket.TextMessage, messageBytes)
			if err != nil {
				return
			}
		}
	}
}

func (h *Handler) MessageHandler(message models.Message) {
	fmt.Printf("%v : %v \n", message.UserNickname, message.Content)

}
func (h *Handler) sendLastMessages(ctx context.Context, connection *websocket.Conn) {
	select {
	case <-ctx.Done():
		return
	default:
		messages, err := h.Services.GetAmountMessage(ctx, 10)
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

}
func (h *Handler) RegisterHandlers() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/chat", h.Echo).Methods("GET")
	return router
}
