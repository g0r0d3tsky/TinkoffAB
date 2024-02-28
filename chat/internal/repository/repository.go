package repository

import (
	"2024-spring-ab-go-hw-1-template-g0r0d3tsky/chat/internal/config"
	"2024-spring-ab-go-hw-1-template-g0r0d3tsky/chat/internal/domain"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type rw struct {
	store *pgxpool.Pool
}

// TODO: mock

type Message interface {
	CreateMessage(ctx context.Context, duty *domain.Message) error
	GetAmountMessage(ctx context.Context, amount int) ([]*domain.Message, error)
}

// go:generate mockery --name ServiceRepository
type ServiceRepository interface {
	Message
}

func New(dbPool *pgxpool.Pool) ServiceRepository {
	return rw{
		store: dbPool,
	}
}
func Connect(c *config.Config) (*pgxpool.Pool, error) {
	connectionString := c.PostgresDSN()

	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx pool config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return pool, nil
}
