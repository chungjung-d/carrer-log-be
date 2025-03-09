package types

import (
	"os"

	"github.com/sashabaranov/go-openai"
)

// CompletionMessage는 ChatGPT API에 보낼 메시지의 기본 인터페이스입니다
type CompletionMessage interface {
	ToContent() (string, error)
}

// ChatGPTConfig는 ChatGPT 서비스 설정을 위한 구조체입니다
type ChatGPTConfig struct {
	APIKey string
	Model  string
}

// DefaultConfig는 기본 설정을 반환합니다
func DefaultConfig() *ChatGPTConfig {
	return &ChatGPTConfig{
		APIKey: os.Getenv("OPENAI_API_KEY"),
		Model:  openai.GPT3Dot5Turbo,
	}
}
