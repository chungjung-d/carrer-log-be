package chat

import (
	appErrors "career-log-be/errors"
	"career-log-be/models/note/chat"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type GetChatResponse struct {
	ID        string        `json:"id"`
	UserID    string        `json:"userId"`
	Title     string        `json:"title"`
	ChatData  chat.ChatData `json:"chatData"`
	CreatedAt string        `json:"createdAt"`
	UpdatedAt string        `json:"updatedAt"`
}

func HandleGetChat(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	chatID := c.Params("id")
	userID := c.Locals("user_id").(string)

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

	resp := GetChatResponse{
		ID:        chatSet.ID,
		UserID:    chatSet.UserID,
		Title:     chatSet.Title,
		ChatData:  chatSet.ChatData,
		CreatedAt: chatSet.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: chatSet.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    resp,
	})
}
