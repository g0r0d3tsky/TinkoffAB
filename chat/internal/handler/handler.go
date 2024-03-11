package handler

import (
	"2024-spring-ab-go-hw-1-template-g0r0d3tsky/chat/internal/domain"
	"2024-spring-ab-go-hw-1-template-g0r0d3tsky/chat/internal/handler/models"
	"context"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type MessageHandler interface {
	Echo(w http.ResponseWriter, r *http.Request)
}

type MessageUseCase interface {
	CreateMessage(ctx context.Context, message *domain.Message) error
	GetAmountMessage(ctx context.Context, amount int) ([]*domain.Message, error)
}

type ChatHandler interface {
	MessageHandler
	MessageHandler(message models.Message)
	RegisterHandlers() http.Handler
}

type Handler struct {
	Services MessageUseCase
	Clients  map[*websocket.Conn]struct{}
	mutex    sync.Mutex
}
