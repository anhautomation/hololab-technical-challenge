package domain

type Bot struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Job       string `json:"job"`
	Bio       string `json:"bio"`
	Style     string `json:"style"`
	Knowledge string `json:"knowledge"`
	CreatedAt string `json:"created_at"`
}
