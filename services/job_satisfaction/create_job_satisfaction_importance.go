package job_satisfaction

import (
	appErrors "career-log-be/errors"
	job_satisfaction "career-log-be/models/job_satisfaction"
	"career-log-be/utils/response"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CreateJobSatisfactionImportanceInput struct {
	Workload          float64 `json:"workload" validate:"required,min=0,max=100"`
	Compensation      float64 `json:"compensation" validate:"required,min=0,max=100"`
	Growth            float64 `json:"growth" validate:"required,min=0,max=100"`
	WorkEnvironment   float64 `json:"workEnvironment" validate:"required,min=0,max=100"`
	WorkRelationships float64 `json:"workRelationships" validate:"required,min=0,max=100"`
	WorkValues        float64 `json:"workValues" validate:"required,min=0,max=100"`
}

func HandleCreateJobSatisfactionImportance() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := c.Locals("db").(*gorm.DB)
		userID := c.Locals("userID").(string)
		input := new(CreateJobSatisfactionImportanceInput)

		if err := c.BodyParser(input); err != nil {
			return appErrors.NewBadRequestError(
				appErrors.ErrorCodeInvalidInput,
				"Invalid request body",
			)
		}
		// 입력값 검증
		validate := validator.New()
		if err := validate.Struct(input); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			return appErrors.NewValidationError(
				appErrors.ErrorCodeInvalidInput,
				"Validation failed",
				validationErrors.Error(),
			)
		}

		// 이미 직무 만족도 중요도가 존재하는지 확인
		var existingImportance job_satisfaction.UserJobSatisfactionImportance
		result := db.Where("user_id = ?", userID).First(&existingImportance)
		if result.Error == nil {
			return appErrors.NewBadRequestError(
				appErrors.ErrorCodeResourceExists,
				"Job satisfaction importance already exists",
			)
		} else if result.Error != gorm.ErrRecordNotFound {
			return appErrors.NewInternalError(
				appErrors.ErrorCodeDatabaseError,
				"Failed to query database",
				result.Error,
			)
		}

		// 새 직무 만족도 중요도 생성
		jobSatisfactionImportance := job_satisfaction.UserJobSatisfactionImportance{
			UserID:            userID,
			Workload:          input.Workload,
			Compensation:      input.Compensation,
			Growth:            input.Growth,
			WorkEnvironment:   input.WorkEnvironment,
			WorkRelationships: input.WorkRelationships,
			WorkValues:        input.WorkValues,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		if err := db.Create(&jobSatisfactionImportance).Error; err != nil {
			return appErrors.NewInternalError(
				appErrors.ErrorCodeDatabaseError,
				"Failed to create job satisfaction importance",
				err,
			)
		}

		return response.Created(c, jobSatisfactionImportance)
	}
}
