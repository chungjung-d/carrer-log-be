package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"career-log-be/utils/chatgpt/types"

	"github.com/sashabaranov/go-openai"
)

// ErrMissingAPIKey는 API 키가 없을 때 발생하는 에러입니다
var ErrMissingAPIKey = errors.New("missing API key")

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

// CompleteChatRequestWithMessage는 특정 타입의 메시지를 처리하여 응답을 반환합니다
func (s *ChatGPTService) CompleteChatRequestWithMessage(ctx context.Context, message types.CompletionMessage) (string, error) {
	content, err := message.ToContent()
	if err != nil {
		return "", err
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
		},
	}

	return s.CompleteChatRequest(ctx, messages)
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

// CompleteChatRequestWithType는 응답을 특정 타입으로 변환하여 반환합니다
func CompleteChatRequestWithType[T any](s *ChatGPTService, ctx context.Context, messages []openai.ChatCompletionMessage) (T, error) {
	var result T

	response, err := s.CompleteChatRequest(ctx, messages)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return result, errors.New("failed to parse response to specified type: " + err.Error())
	}

	return result, nil
}
