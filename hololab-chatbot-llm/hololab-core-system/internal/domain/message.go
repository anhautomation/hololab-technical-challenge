package domain

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	BotID     string `json:"bot_id"`
	Role      Role   `json:"role"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}
