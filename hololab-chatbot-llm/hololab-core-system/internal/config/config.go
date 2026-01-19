package config

import "os"

type Config struct {
	Port       string
	SQLiteFile string

	LLMMode string

	OpenAIBaseURL string
	OpenAIAPIKey  string
	OpenAIModel   string
}

func Load() Config {
	c := Config{
		Port:          getenv("PORT", "3001"),
		SQLiteFile:    getenv("SQLITE_FILE", "./data.sqlite"),
		LLMMode:       getenv("LLM_MODE", "mock"),
		OpenAIBaseURL: getenv("OPENAI_BASE_URL", "https://api.openai.com"),
		OpenAIAPIKey:  os.Getenv("OPENAI_API_KEY"),
		OpenAIModel:   getenv("OPENAI_MODEL", "gpt-4o-mini"),
	}
	return c
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
