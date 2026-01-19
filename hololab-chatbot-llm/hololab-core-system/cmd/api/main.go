package main

import (
	"context"
	httpadapter "hololab-core-system/internal/adapters/http"
	"hololab-core-system/internal/adapters/llm"
	"hololab-core-system/internal/adapters/sqlite"
	"hololab-core-system/internal/app/bots"
	"hololab-core-system/internal/app/chat"
	"hololab-core-system/internal/config"
	"hololab-core-system/internal/ports"
	"log"
	"time"

	"github.com/joho/godotenv"
)

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

type timeID struct{}

func (timeID) NewID() string { return time.Now().Format("20060102150405.000000000") }

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	db, err := sqlite.Open(cfg.SQLiteFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	botRepo := sqlite.BotRepo{DB: db}
	msgRepo := sqlite.MessageRepo{DB: db}

	var llmProvider ports.LLMProvider
	if cfg.LLMMode == "openai_compat" {
		llmProvider = llm.OpenAICompatProvider{
			BaseURL: cfg.OpenAIBaseURL,
			APIKey:  cfg.OpenAIAPIKey,
			Model:   cfg.OpenAIModel,
		}
	} else {
		llmProvider = llm.MockProvider{}
	}

	botsSvc := bots.Service{
		Repo:  botRepo,
		Clock: realClock{},
		IDGen: timeID{},
	}
	chatSvc := chat.Service{
		Bots:     botRepo,
		Messages: msgRepo,
		LLM:      llmProvider,
		Clock:    realClock{},
	}

	r := httpadapter.NewRouter(httpadapter.RouterDeps{
		Bots: botsSvc,
		Chat: chatSvc,
	})

	log.Printf("[hololab-core-system] listening on :%s (db=%s, mode=%s)", cfg.Port, cfg.SQLiteFile, cfg.LLMMode)
	_ = r.Run(":" + cfg.Port)

	_ = context.Background()
}
