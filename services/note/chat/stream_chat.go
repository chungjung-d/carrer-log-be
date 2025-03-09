package chat

import (
	"bufio"
	appErrors "career-log-be/errors"
	"career-log-be/models/note/chat"
	"career-log-be/models/note/chat/enums"
	"career-log-be/utils/chatgpt"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sashabaranov/go-openai"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
)

type StreamChatRequest struct {
	Message string `json:"message"`
}

func HandleStreamChat(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	chatGPTService := c.Locals("chatgpt").(chatgpt.Service)
	userID := c.Locals("user_id").(string)
	chatID := c.Params("id")

	var req StreamChatRequest
	if err := c.BodyParser(&req); err != nil {
		return appErrors.NewBadRequestError(
			"Invalid request body",
			err.Error(),
		)
	}

	// ChatSet 조회
	var chatSet chat.ChatSet
	result := db.Where("id = ? AND user_id = ?", chatID, userID).First(&chatSet)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return appErrors.NewNotFoundError(
				appErrors.ErrorCodeResourceNotFound,
				"Chat not found",
			)
		}
		return appErrors.NewInternalError(
			appErrors.ErrorCodeDatabaseError,
			"Failed to retrieve chat",
			result.Error,
		)
	}

	// 사용자 메시지 추가
	chatSet.ChatData.AddMessage(enums.UserRole, req.Message)

	// ChatGPT 메시지 준비
	messages := make([]openai.ChatCompletionMessage, len(chatSet.ChatData.Messages))
	for i, msg := range chatSet.ChatData.Messages {
		messages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role.String(),
			Content: msg.Content,
		}
	}

	// SSE 설정
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	// 스트리밍 응답 시작
	responseChan, errChan := chatGPTService.StreamChatRequest(context.Background(), messages)

	// 임시 버퍼에 응답 저장
	var fullResponse string

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		defer w.Flush()

		for {
			select {
			case chunk, ok := <-responseChan:
				if !ok {
					// 스트림 종료
					chatSet.ChatData.AddMessage(enums.AssistantRole, fullResponse)

					// DB 업데이트
					if err := db.Save(&chatSet).Error; err != nil {
						fmt.Printf("Failed to save chat: %v\n", err)
					}
					return
				}

				fullResponse += chunk

				// SSE 형식으로 데이터 전송
				data, _ := json.Marshal(fiber.Map{
					"content": chunk,
				})
				w.Write([]byte(fmt.Sprintf("data: %s\n\n", data)))
				w.Flush()

			case err, ok := <-errChan:
				if !ok {
					return
				}
				fmt.Printf("Streaming error: %v\n", err)
				return
			}
		}
	}))

	return nil
}
