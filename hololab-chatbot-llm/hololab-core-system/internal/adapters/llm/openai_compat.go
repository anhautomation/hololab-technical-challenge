package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"hololab-core-system/internal/ports"
	"net/http"
	"strings"
	"time"
)

type OpenAICompatProvider struct {
	BaseURL string
	APIKey  string
	Model   string
}

type reqBody struct {
	Model       string          `json:"model"`
	Messages    []ports.ChatMsg `json:"messages"`
	Temperature float64         `json:"temperature"`
}
type respBody struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (p OpenAICompatProvider) Chat(ctx context.Context, systemPrompt string, history []ports.ChatMsg, userText string) (string, error) {
	if strings.TrimSpace(p.APIKey) == "" {
		return "", errors.New("missing api key")
	}
	base := strings.TrimRight(p.BaseURL, "/")
	if base == "" {
		base = "https://api.openai.com"
	}
	model := p.Model
	if model == "" {
		model = "gpt-4o-mini"
	}

	msgs := make([]ports.ChatMsg, 0, 2+len(history))
	msgs = append(msgs, ports.ChatMsg{Role: "system", Content: systemPrompt})
	msgs = append(msgs, history...)
	msgs = append(msgs, ports.ChatMsg{Role: "user", Content: userText})

	b, _ := json.Marshal(reqBody{Model: model, Messages: msgs, Temperature: 0.7})

	req, _ := http.NewRequestWithContext(ctx, "POST", base+"/v1/chat/completions", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.APIKey)

	client := &http.Client{Timeout: 40 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", errors.New("llm error status: " + res.Status)
	}

	var out respBody
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", errors.New("empty reply")
	}
	content := strings.TrimSpace(out.Choices[0].Message.Content)
	if content == "" {
		return "(empty reply)", nil
	}
	return content, nil
}
