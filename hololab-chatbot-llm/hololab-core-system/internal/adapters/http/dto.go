package httpadapter

type CreateBotRequest struct {
	Name      string `json:"name"`
	Job       string `json:"job"`
	Bio       string `json:"bio"`
	Style     string `json:"style"`
	Knowledge string `json:"knowledge"`
}

type SendMessageRequest struct {
	Message string `json:"message"`
}
