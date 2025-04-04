package chat

import (
	appErrors "career-log-be/errors"
	"career-log-be/models/note/chat"
	"career-log-be/utils/response"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ListPreChatsResponse struct {
	PreChats []chat.PreChat `json:"pre_chats"`
}

func HandleListPreChats(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)

	var preChats []chat.PreChat
	if err := db.Order("created_at desc").Find(&preChats).Error; err != nil {
		return appErrors.NewInternalError(
			appErrors.ErrorCodeDatabaseError,
			"Failed to retrieve pre-chats",
			err,
		)
	}

	resp := ListPreChatsResponse{
		PreChats: preChats,
	}

	return response.Success(c, resp)
}
