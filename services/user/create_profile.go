package user

import (
	appErrors "career-log-be/errors"
	user "career-log-be/models/user"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var validate = validator.New()

type CreateUserProfileInput struct {
	Name         string `json:"name" validate:"required"`
	Nickname     string `json:"nickname" validate:"required"`
	Organization string `json:"organization" validate:"required"`
}

func HandleCreateUserProfile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := c.Locals("db").(*gorm.DB)
		userID := c.Locals("userID").(string)
		input := new(CreateUserProfileInput)

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

		// 이미 프로필이 존재하는지 확인
		var existingProfile user.UserProfile
		result := db.Where("id = ?", userID).First(&existingProfile)
		if result.Error == nil {
			return appErrors.NewBadRequestError(
				appErrors.ErrorCodeResourceExists,
				"User profile already exists",
			)
		} else if result.Error != gorm.ErrRecordNotFound {
			return appErrors.NewInternalError(
				appErrors.ErrorCodeDatabaseError,
				"Failed to query database",
				result.Error,
			)
		}

		// 새 프로필 생성
		userProfile := user.UserProfile{
			ID:           userID,
			Name:         input.Name,
			Nickname:     input.Nickname,
			Organization: input.Organization,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := db.Create(&userProfile).Error; err != nil {
			return appErrors.NewInternalError(
				appErrors.ErrorCodeDatabaseError,
				"Failed to create user profile",
				err,
			)
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"data":    userProfile,
		})
	}
}
