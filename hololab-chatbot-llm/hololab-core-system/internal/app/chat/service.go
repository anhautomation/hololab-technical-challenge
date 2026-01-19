package chat

import (
	"context"
	"hololab-core-system/internal/app/prompt"
	"hololab-core-system/internal/domain"
	"hololab-core-system/internal/ports"
	"strings"
)

type Service struct {
	Bots     ports.BotRepository
	Messages ports.MessageRepository
	LLM      ports.LLMProvider
	Clock    ports.Clock
}

func (s Service) History(ctx context.Context, botID string) ([]domain.Message, error) {
	_, err := s.Bots.Get(ctx, botID)
	if err != nil {
		return nil, err
	}
	return s.Messages.ListByBot(ctx, botID)
}

func (s Service) Reset(ctx context.Context, botID string) error {
	_, err := s.Bots.Get(ctx, botID)
	if err != nil {
		return err
	}
	return s.Messages.ResetByBot(ctx, botID)
}

func (s Service) Send(ctx context.Context, botID string, userText string) (reply string, err error) {
	bot, err := s.Bots.Get(ctx, botID)
	if err != nil {
		return "", err
	}

	userText = strings.TrimSpace(userText)
	if userText == "" {
		return "", domain.ErrInvalidInput
	}

	history, err := s.Messages.ListByBot(ctx, botID)
	if err != nil {
		return "", err
	}

	_ = s.Messages.Append(ctx, domain.Message{
		BotID:     botID,
		Role:      domain.RoleUser,
		Content:   userText,
		CreatedAt: s.Clock.Now().Format("2006-01-02T15:04:05Z07:00"),
	})

	var llmHist []ports.ChatMsg
	if len(history) > 12 {
		history = history[len(history)-12:]
	}
	for _, m := range history {
		role := "user"
		if m.Role == domain.RoleAssistant {
			role = "assistant"
		}
		llmHist = append(llmHist, ports.ChatMsg{Role: role, Content: m.Content})
	}

	system := prompt.BuildSystemPrompt(bot)

	reply, err = s.LLM.Chat(ctx, system, llmHist, userText)
	if err != nil {
		return "", domain.ErrLLMUnavailable
	}

	_ = s.Messages.Append(ctx, domain.Message{
		BotID:     botID,
		Role:      domain.RoleAssistant,
		Content:   reply,
		CreatedAt: s.Clock.Now().Format("2006-01-02T15:04:05Z07:00"),
	})

	return reply, nil
}
