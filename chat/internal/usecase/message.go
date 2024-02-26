package usecase

import (
	"2024-spring-ab-go-hw-1-template-g0r0d3tsky/chat/internal/domain"
	"2024-spring-ab-go-hw-1-template-g0r0d3tsky/chat/internal/repository"
	"context"
)

type UC struct {
	Repo repository.ServiceRepository
}
type MessageUseCase interface {
	CreateMessage(ctx context.Context, message *domain.Message) error
	GetAmountMessage(ctx context.Context, amount int) ([]*domain.Message, error)
}

func (uc *UC) CreateMessage(ctx context.Context, message *domain.Message) error {
	return uc.Repo.CreateMessage(ctx, message)
}

func (uc *UC) GetAmountMessage(ctx context.Context, amount int) ([]*domain.Message, error) {
	return uc.Repo.GetAmountMessage(ctx, amount)
}
