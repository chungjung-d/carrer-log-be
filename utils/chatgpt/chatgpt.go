// Package chatgpt provides a ChatGPT client implementation
package chatgpt

import (
	"career-log-be/utils/chatgpt/builder"
	"career-log-be/utils/chatgpt/service"
	"career-log-be/utils/chatgpt/types"
)

type (
	ChatGPTConfig = types.ChatGPTConfig
)

// Service는 ChatGPT 서비스의 인터페이스입니다
type Service = service.ChatGPTService

// NewChatGPTBuilder는 새로운 ChatGPT 빌더를 생성합니다
func NewChatGPTBuilder() *builder.ChatGPTBuilder {
	return builder.NewChatGPTBuilder()
}
