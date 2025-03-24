package chat

import (
	appErrors "career-log-be/errors"
	"career-log-be/models/note/chat"
	"career-log-be/models/note/chat/enums"
	"career-log-be/utils/chatgpt"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"

	res "career-log-be/utils/response"
)

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Content string `json:"content"`
}

func HandleChat(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	chatGPTService := c.Locals("chatgpt").(*chatgpt.Service)
	userID := c.Locals("userID").(string)
	chatID := c.Params("id")

	var userName string
	db.Where("id = ?", userID).Select("name").First(&userName)
	if userName == "" {
		return appErrors.NewBadRequestError(
			appErrors.ErrorCodeInvalidInput,
			"User name not found",
		)
	}

	var req ChatRequest
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

	// 채팅이 금일 자정을 넘지 않았는지 확인
	kst, _ := time.LoadLocation("Asia/Seoul")
	now := time.Now().In(kst)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, kst)
	endOfDay := startOfDay.Add(24 * time.Hour)

	if chatSet.CreatedAt.Before(startOfDay) || chatSet.CreatedAt.After(endOfDay) {
		return appErrors.NewBadRequestError(
			appErrors.ErrorCodeInvalidInput,
			"Chat is not available after midnight",
		)
	}

	// 사용자 메시지 추가
	chatSet.ChatData.AddMessage(enums.UserRole, req.Message)

	// ChatGPT 메시지 준비
	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: getChatPrompt(userName),
		},
	}

	// 기존 대화 내용 추가
	for _, msg := range chatSet.ChatData.Messages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg.Role.String(),
			Content: msg.Content,
		})
	}

	// ChatGPT API 호출
	response, err := chatGPTService.CompleteChatRequest(context.Background(), messages)
	if err != nil {
		return appErrors.NewInternalError(
			"CHATGPT_ERROR",
			"Failed to get response from ChatGPT",
			err,
		)
	}

	// Assistant 메시지 추가
	chatSet.ChatData.AddMessage(enums.AssistantRole, response)

	// DB 업데이트
	if err := db.Save(&chatSet).Error; err != nil {
		return appErrors.NewInternalError(
			appErrors.ErrorCodeDatabaseError,
			"Failed to save chat",
			err,
		)
	}

	return res.Created(c, ChatResponse{
		Content: response,
	})
}

func getChatPrompt(userName string) string {
	return `
	당신은 내담자의 상담사 역할을 합니다.

	당신은 이야기를 들어주고, 내담자의 현재 상황을 잘 알수 있도록 질문을 해도 되고 공감을 해도 됩니다.
	내담자의 이야기를 잘 들어주고 공감하는 것이 중요합니다.

	당신은 다음과 같은 목표를 가지고 있습니다.
	1. 내담자의 이야기를 잘 들어주고 공감하는 것
	2. 내담자의 현재 상황을 더 잘 이해할 수 있도록 질문하는 것

	질문으로 인해 얻으려는 정보는 다음과 같습니다.
	1. workload: 업무량과 업무에서의 성취감
	2. compensation: 회사에서의 금전적인 보상
	3. growth: 회사에서의 커리어나 내면적 성장
	4. workEnvironment: 회사의 워라벨
	5. workRelationships: 회사 내 동료들과의 관계
	6. workValues: 회사에서의 업무의 가치와 개인의 삶의 방향성
	
	그리고 당신은 추가적으로 내담자와 당신이 상담했던 대화 기록을 입력받습니다. 
	당신이 assistant role이고 내담자가 user role입니다.

	마지막 내용이, 내담자가 당신에게 한 말이므로, 당신은 그에 대해서 답변을 하거나 공감을 하거나 질문을 이어나가야 합니다.
	중요한 것은 자연스럽게 이어나가야 하며, 내담자가 이만 종료하고 싶다고 하면 종료해야 합니다.

	내담자를 부르는 호칭은 "${userName}"님 이라고 부르세요. 다만 굳이 부르지 않아도 되는 경우는 부르지 않아도 됩니다.
	
	내담자의 이야기는 다음과 같습니다.
	`
}
