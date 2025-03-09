package service

import (
	"context"
	"errors"
	"io"

	"career-log-be/utils/chatgpt/types"

	"github.com/sashabaranov/go-openai"
)

// ErrMissingAPIKey는 API 키가 없을 때 발생하는 에러입니다
var ErrMissingAPIKey = errors.New("missing API key")

// ChatGPTResponse는 ChatGPT로부터 받을 수 있는 구조화된 응답의 예시입니다
type ChatGPTResponse struct {
	Success    bool     `json:"success"`
	Message    string   `json:"message"`
	Confidence float64  `json:"confidence"`
	Categories []string `json:"categories"`
}

// ChatGPTService는 ChatGPT API 서비스를 제공합니다
type ChatGPTService struct {
	client *openai.Client
	config *types.ChatGPTConfig
}

// NewChatGPTService는 새로운 ChatGPT 서비스를 생성합니다
func NewChatGPTService(client *openai.Client, config *types.ChatGPTConfig) *ChatGPTService {
	return &ChatGPTService{
		client: client,
		config: config,
	}
}

// CompleteChatRequest는 일반적인 채팅 완료 요청을 처리합니다
func (s *ChatGPTService) CompleteChatRequest(ctx context.Context, messages []openai.ChatCompletionMessage) (string, error) {
	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    s.config.Model,
			Messages: messages,
		},
	)

	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no response choices available")
	}

	return resp.Choices[0].Message.Content, nil
}

// StreamChatRequest는 스트리밍 방식으로 채팅 완료 요청을 처리합니다
func (s *ChatGPTService) StreamChatRequest(ctx context.Context, messages []openai.ChatCompletionMessage) (chan string, chan error) {
	responseChan := make(chan string)
	errChan := make(chan error)

	go func() {
		defer close(responseChan)
		defer close(errChan)

		stream, err := s.client.CreateChatCompletionStream(
			ctx,
			openai.ChatCompletionRequest{
				Model:    s.config.Model,
				Messages: messages,
			},
		)
		if err != nil {
			errChan <- err
			return
		}
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return
			}

			if err != nil {
				errChan <- err
				return
			}

			if len(response.Choices) > 0 {
				responseChan <- response.Choices[0].Delta.Content
			}
		}
	}()

	return responseChan, errChan
}

// CompleteChatRequestWithMessage는 특정 타입의 메시지를 처리하여 응답을 반환합니다
func CompleteChatRequestWithMessage[T any](ctx context.Context, s *ChatGPTService, messages []openai.ChatCompletionMessage, parser func(string) (T, error)) (T, error) {
	response, err := s.CompleteChatRequest(ctx, messages)
	if err != nil {
		var zero T
		return zero, err
	}

	return parser(response)
}

// 구조체 사용 예시:
/*
func ExampleStructureUsage() {
	parser := func(response string) (ChatGPTResponse, error) {
		var result ChatGPTResponse
		if err := json.Unmarshal([]byte(response), &result); err != nil {
			return ChatGPTResponse{}, err
		}
		return result, nil
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    "user",
			Content: "이 텍스트가 긍정적인가요? 응답을 JSON 형식으로 주세요: '오늘은 정말 좋은 날이에요!'",
		},
	}

	result, err := CompleteChatRequestWithMessage(context.Background(), chatGPTService, messages, parser)
	if err != nil {
		// 에러 처리
		return
	}

	// result.Success, result.Message 등의 필드 사용 가능
}
*/
