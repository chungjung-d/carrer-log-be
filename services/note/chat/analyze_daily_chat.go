package chat

import (
	"career-log-be/services/note/chat/scheduler"
	"career-log-be/utils/response"

	appErrors "career-log-be/errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// HandleAnalyzeDailyChat 일일 채팅 분석을 수동으로 실행하는 핸들러
func HandleAnalyzeDailyChat(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)

	chatScheduler, err := scheduler.NewChatAnalyzeScheduler(db)
	if err != nil {
		return appErrors.NewInternalError(
			appErrors.ErrorCodeInternalError,
			"Failed to create chat analyzer",
			err,
		)
	}

	chatScheduler.AnalyzeDailyChat()

	return response.Accepted(c, "Daily chat analysis has been executed")
}
