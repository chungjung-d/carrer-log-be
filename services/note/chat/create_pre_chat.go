package chat

import (
	appErrors "career-log-be/errors"
	"career-log-be/models/note/chat"
	"career-log-be/utils/response"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CreatePreChatRequest struct {
	Content string `json:"content" validate:"required"`
}

type CreatePreChatResponse struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func HandleCreatePreChat(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)

	var req CreatePreChatRequest
	if err := c.BodyParser(&req); err != nil {
		return appErrors.NewBadRequestError(
			appErrors.ErrorCodeInvalidInput,
			"Invalid request body",
		)
	}

	// 입력값 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return appErrors.NewValidationError(
			appErrors.ErrorCodeInvalidInput,
			"Validation failed",
			validationErrors.Error(),
		)
	}

	// 새로운 PreChat 생성
	kst, _ := time.LoadLocation("Asia/Seoul")
	now := time.Now().In(kst)

	preChat := chat.PreChat{
		Content:   req.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// DB에 저장
	if err := db.Create(&preChat).Error; err != nil {
		return appErrors.NewInternalError(
			appErrors.ErrorCodeDatabaseError,
			"Failed to create pre-chat",
			err,
		)
	}

	resp := CreatePreChatResponse{
		ID:        preChat.ID,
		Content:   preChat.Content,
		CreatedAt: preChat.CreatedAt,
	}

	return response.Created(c, resp)
}
