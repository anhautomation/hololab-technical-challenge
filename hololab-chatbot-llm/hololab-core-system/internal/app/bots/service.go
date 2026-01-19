package bots

import (
	"context"
	"hololab-core-system/internal/domain"
	"hololab-core-system/internal/ports"
	"strings"
)

type Service struct {
	Repo  ports.BotRepository
	Clock ports.Clock
	IDGen ports.IDGenerator
}

type CreateInput struct {
	Name      string
	Job       string
	Bio       string
	Style     string
	Knowledge string
}

func (s Service) List(ctx context.Context) ([]domain.Bot, error) {
	return s.Repo.List(ctx)
}

func (s Service) Get(ctx context.Context, id string) (domain.Bot, error) {
	return s.Repo.Get(ctx, id)
}

func (s Service) Create(ctx context.Context, in CreateInput) (domain.Bot, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.Job = strings.TrimSpace(in.Job)
	in.Bio = strings.TrimSpace(in.Bio)
	in.Style = strings.TrimSpace(in.Style)
	in.Knowledge = strings.TrimSpace(in.Knowledge)

	if in.Name == "" || in.Job == "" || in.Bio == "" || in.Style == "" {
		return domain.Bot{}, domain.ErrInvalidInput
	}

	b := domain.Bot{
		ID:        s.IDGen.NewID(),
		Name:      in.Name,
		Job:       in.Job,
		Bio:       in.Bio,
		Style:     in.Style,
		Knowledge: in.Knowledge,
		CreatedAt: s.Clock.Now().Format("2006-01-02T15:04:05Z07:00"),
	}

	if err := s.Repo.Create(ctx, b); err != nil {
		return domain.Bot{}, err
	}
	return b, nil
}

func (s Service) Delete(ctx context.Context, id string) (int64, error) {
	return s.Repo.Delete(ctx, id)
}
