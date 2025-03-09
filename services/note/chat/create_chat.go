package chat

import (
	appErrors "career-log-be/errors"
	"career-log-be/models/note/chat"
	"career-log-be/models/note/chat/enums"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CreateChatRequest struct {
	PreChatID string `json:"pre_chat_id"`
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

	return c.JSON(fiber.Map{
		"id": chatSet.ID,
	})
}
