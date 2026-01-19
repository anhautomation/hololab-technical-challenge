package ports

import (
	"context"
	"hololab-core-system/internal/domain"
)

type MessageRepository interface {
	ListByBot(ctx context.Context, botID string) ([]domain.Message, error)
	Append(ctx context.Context, msg domain.Message) error
	ResetByBot(ctx context.Context, botID string) error
}
