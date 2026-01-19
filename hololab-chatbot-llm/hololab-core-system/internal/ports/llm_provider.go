package ports

import "context"

type ChatMsg struct {
	Role    string
	Content string
}

type LLMProvider interface {
	Chat(ctx context.Context, systemPrompt string, history []ChatMsg, userText string) (string, error)
}
