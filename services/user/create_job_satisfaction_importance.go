package user

import (
	appErrors "career-log-be/errors"
	"career-log-be/models"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CreateJobSatisfactionImportanceInput struct {
	Workload          int `json:"workload" validate:"required,min=0,max=100"`
	Compensation      int `json:"compensation" validate:"required,min=0,max=100"`
	Growth            int `json:"growth" validate:"required,min=0,max=100"`
	WorkEnvironment   int `json:"workEnvironment" validate:"required,min=0,max=100"`
	WorkRelationships int `json:"workRelationships" validate:"required,min=0,max=100"`
	WorkValues        int `json:"workValues" validate:"required,min=0,max=100"`
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
		if err := validate.Struct(input); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			return appErrors.NewValidationError(
				appErrors.ErrorCodeInvalidInput,
				"Validation failed",
				validationErrors.Error(),
			)
		}

		// 이미 직무 만족도 중요도가 존재하는지 확인
		var existingImportance models.UserJobSatisfactionImportance
		result := db.Where("id = ?", userID).First(&existingImportance)
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
		jobSatisfactionImportance := models.UserJobSatisfactionImportance{
			ID:                userID,
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

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"data":    jobSatisfactionImportance,
		})
	}
}
