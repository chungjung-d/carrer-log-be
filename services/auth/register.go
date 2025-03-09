package auth

import (
	appErrors "career-log-be/errors"
	"career-log-be/models/user"
	"career-log-be/utils/response"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var validate = validator.New()

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// if err := validate.Struct(input); err != nil {
// 	validationErrors := err.(validator.ValidationErrors)
// 	return appErrors.NewValidationError(
// 		appErrors.ErrorCodeInvalidInput,
// 		"Validation failed",
// 		validationErrors.Error(),
// 	)
// }

func HandleRegister() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// DB 인스턴스 가져오기
		db := c.Locals("db").(*gorm.DB)

		input := new(RegisterInput)

		if err := c.BodyParser(input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input",
			})
		}

		// Validate the input
		if err := validate.Struct(input); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			return appErrors.NewValidationError(
				appErrors.ErrorCodeInvalidInput,
				"Validation failed",
				validationErrors.Error(),
			)
		}

		// 이메일 중복 체크
		var existingUser user.User
		result := db.Where("email = ?", input.Email).First(&existingUser)
		if result.Error == nil {
			return appErrors.NewConflictError(
				appErrors.ErrorCodeInvalidInput,
				"Email already exists",
			)
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return appErrors.NewInternalError(
				appErrors.ErrorCodeDatabaseError,
				"Database error",
				result.Error,
			)
		}

		// 비밀번호 해싱
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not hash password",
			})
		}

		// 새 사용자 생성
		user := &user.User{
			Email:    input.Email,
			Password: string(hashedPassword),
		}

		if err := db.Create(user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not create user",
			})
		}

		// 응답에서 비밀번호 제외
		resp := RegisterResponse{
			ID:    user.ID,
			Email: user.Email,
		}

		return response.Created(c, resp)
	}
}
