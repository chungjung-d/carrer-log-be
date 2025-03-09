package chat

import (
	appErrors "career-log-be/errors"
	"career-log-be/models/note/chat"
	"career-log-be/models/note/chat/enums"
	"career-log-be/utils/response"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CreateChatRequest struct {
	PreChatID string `json:"pre_chat_id"`
}

type CreateChatResponse struct {
	ID string `json:"id"`
}

func HandleCreateChat(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	userID := c.Locals("user_id").(string)

	var req CreateChatRequest
	if err := c.BodyParser(&req); err != nil {
		return appErrors.NewBadRequestError(
			"Invalid request body",
			err.Error(),
		)
	}

	// 오늘 생성한 채팅이 있는지 확인
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var existingChat chat.ChatSet
	if err := db.Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, startOfDay, endOfDay).First(&existingChat).Error; err == nil {
		return appErrors.NewBadRequestError(
			appErrors.ErrorCodeInvalidInput,
			"Daily chat limit exceeded",
		)
	} else if err != gorm.ErrRecordNotFound {
		return appErrors.NewInternalError(
			appErrors.ErrorCodeDatabaseError,
			"Failed to check existing chat",
			err,
		)
	}

	// PreChat 조회
	var preChat chat.PreChat
	if err := db.Where("id = ?", req.PreChatID).First(&preChat).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.NewNotFoundError(
				appErrors.ErrorCodeResourceNotFound,
				"PreChat not found",
			)
		}
		return appErrors.NewInternalError(
			appErrors.ErrorCodeDatabaseError,
			"Failed to retrieve pre-chat",
			err,
		)
	}

	// 새로운 ChatSet 생성
	chatSet := chat.ChatSet{
		UserID: userID,
		ChatData: chat.ChatData{
			Messages: []chat.Message{
				{
					Role:    enums.AssistantRole,
					Content: preChat.Content,
				},
			},
		},
	}

	// DB에 저장
	if err := db.Create(&chatSet).Error; err != nil {
		return appErrors.NewInternalError(
			appErrors.ErrorCodeDatabaseError,
			"Failed to create chat",
			err,
		)
	}

	resp := CreateChatResponse{
		ID: chatSet.ID,
	}

	return response.Created(c, resp)
}
