package builder

import (
	"career-log-be/utils/chatgpt/service"
	"career-log-be/utils/chatgpt/types"

	"github.com/sashabaranov/go-openai"
)

// ChatGPTBuilder는 ChatGPT 서비스 빌더입니다
type ChatGPTBuilder struct {
	config *types.ChatGPTConfig
}

// NewChatGPTBuilder는 새로운 ChatGPT 빌더를 생성합니다
func NewChatGPTBuilder() *ChatGPTBuilder {
	return &ChatGPTBuilder{
		config: types.DefaultConfig(),
	}
}

// WithAPIKey는 API 키를 설정합니다
func (b *ChatGPTBuilder) WithAPIKey(apiKey string) *ChatGPTBuilder {
	b.config.APIKey = apiKey
	return b
}

// WithModel은 모델을 설정합니다
func (b *ChatGPTBuilder) WithModel(model string) *ChatGPTBuilder {
	b.config.Model = model
	return b
}

// Build는 ChatGPT 서비스를 생성합니다
func (b *ChatGPTBuilder) Build() (*service.ChatGPTService, error) {
	if b.config.APIKey == "" {
		return nil, service.ErrMissingAPIKey
	}

	client := openai.NewClient(b.config.APIKey)
	return service.NewChatGPTService(client, b.config), nil
}
