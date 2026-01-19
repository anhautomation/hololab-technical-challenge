package ports

import (
	"context"
	"hololab-core-system/internal/domain"
)

type BotRepository interface {
	List(ctx context.Context) ([]domain.Bot, error)
	Get(ctx context.Context, id string) (domain.Bot, error)
	Create(ctx context.Context, bot domain.Bot) error
	Delete(ctx context.Context, id string) (deleted int64, err error)
}
