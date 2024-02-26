package repository

import (
	"2024-spring-ab-go-hw-1-template-g0r0d3tsky/chat/internal/domain"
	"context"
	"github.com/google/uuid"
)

func (rw rw) CreateMessage(ctx context.Context, mes *domain.Message) error {
	id := uuid.New()
	mes.ID = id
	_, err := rw.store.Exec(ctx,
		`INSERT INTO "messages" (id, user_nickname, content, time) VALUES($1, $2, $3, $4)`,
		&mes.ID, &mes.UserNickname, &mes.Content, &mes.Time,
	)
	if err != nil {
		return err
	}
	return nil
}

func (rw rw) GetAmountMessage(ctx context.Context, amount int) ([]*domain.Message, error) {
	var messages []*domain.Message

	rows, err := rw.store.Query(
		ctx,
		`SELECT id, user_nickname, content, time FROM messages
				ORDER BY time DESC
				LIMIT $1;`, amount)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		mes := &domain.Message{}
		if err := rows.Scan(&mes.ID, &mes.UserNickname, &mes.Content, &mes.Time); err != nil {
			return nil, err
		}

		messages = append(messages, mes)
	}

	return messages, nil
}
