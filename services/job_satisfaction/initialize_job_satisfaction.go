package job_satisfaction

import (
	appErrors "career-log-be/errors"
	job_satisfaction "career-log-be/models/job_satisfaction"
	enums "career-log-be/models/job_satisfaction/enums"
	satisfaction "career-log-be/services/job_satisfaction/core/event"
	"career-log-be/utils/response"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type InitializeJobSatisfactionInput struct {
	Workload          float64 `json:"workload" validate:"required,min=0,max=100"`
	Compensation      float64 `json:"compensation" validate:"required,min=0,max=100"`
	Growth            float64 `json:"growth" validate:"required,min=0,max=100"`
	WorkEnvironment   float64 `json:"workEnvironment" validate:"required,min=0,max=100"`
	WorkRelationships float64 `json:"workRelationships" validate:"required,min=0,max=100"`
	WorkValues        float64 `json:"workValues" validate:"required,min=0,max=100"`
}

func HandleInitializeJobSatisfaction() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := c.Locals("db").(*gorm.DB)
		userID := c.Locals("userID").(string)
		input := new(InitializeJobSatisfactionInput)

		if err := c.BodyParser(input); err != nil {
			return appErrors.NewBadRequestError(
				appErrors.ErrorCodeInvalidInput,
				"Invalid request body",
			)
		}

		validate := validator.New()
		if err := validate.Struct(input); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			return appErrors.NewValidationError(
				appErrors.ErrorCodeInvalidInput,
				"Validation failed",
				validationErrors.Error(),
			)
		}

		// 이미 직무 만족도가 존재하는지 확인
		var existingSatisfaction job_satisfaction.UserJobSatisfaction
		result := db.Where("user_id = ?", userID).First(&existingSatisfaction)
		if result.Error == nil {
			return appErrors.NewBadRequestError(
				appErrors.ErrorCodeResourceExists,
				"Job satisfaction already exists",
			)
		} else if result.Error != gorm.ErrRecordNotFound {
			return appErrors.NewInternalError(
				appErrors.ErrorCodeDatabaseError,
				"Failed to query database",
				result.Error,
			)
		}

		// 이벤트 생성
		updateEvent := &job_satisfaction.JobSatisfactionUpdateEvent{
			UserID:            userID,
			Workload:          input.Workload,
			Compensation:      input.Compensation,
			Growth:            input.Growth,
			WorkEnvironment:   input.WorkEnvironment,
			WorkRelationships: input.WorkRelationships,
			WorkValues:        input.WorkValues,
			EventType:         enums.InitEvent,
			CreatedAt:         time.Now(),
		}

		// 비동기적으로 이벤트 처리
		satisfaction.PublishJobSatisfactionUpdateEvent(db, updateEvent)

		return response.Accepted(c, "Job satisfaction initialization request has been accepted")
	}
}
