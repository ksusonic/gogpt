package models

type YaARTRequest struct {
	UserName string
	Prompt   string

	ChatID    int64
	MessageID int
}

type YaARTResponse struct {
	Image []byte
	Err   GenerationError

	ChatID    int64
	MessageID int
}
