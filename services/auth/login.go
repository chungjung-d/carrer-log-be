package auth

import (
	appErrors "career-log-be/errors"
	"career-log-be/models/user"
	"career-log-be/utils/jwt"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
}

func HandleLogin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := c.Locals("db").(*gorm.DB)
		jwtUtils := c.Locals("jwt").(*jwt.JWTUtils)
		input := new(LoginInput)

		if err := c.BodyParser(input); err != nil {
			return appErrors.NewBadRequestError(
				appErrors.ErrorCodeInvalidInput,
				"Invalid request body",
			)
		}

		// Validate input
		if err := validate.Struct(input); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			return appErrors.NewValidationError(
				appErrors.ErrorCodeInvalidInput,
				"Validation failed",
				validationErrors.Error(),
			)
		}

		// 사용자 찾기
		var user user.User
		result := db.Where("email = ?", input.Email).First(&user)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return appErrors.NewAuthorizationError(
					appErrors.ErrorCodeInvalidCredentials,
					"Invalid email or password",
				)
			}
			return appErrors.NewInternalError(
				appErrors.ErrorCodeDatabaseError,
				"Failed to query database",
				result.Error,
			)
		}

		// 비밀번호 확인
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			return appErrors.NewAuthorizationError(
				appErrors.ErrorCodeInvalidCredentials,
				"Invalid email or password",
			)
		}

		// JWT 토큰 생성
		token, err := jwtUtils.GenerateToken(user.ID, user.Email)
		if err != nil {
			return appErrors.NewInternalError(
				appErrors.ErrorCodeInternalError,
				"Could not generate token",
				err,
			)
		}

		// 응답 생성
		response := LoginResponse{
			Token: token,
		}
		response.User.ID = user.ID
		response.User.Email = user.Email

		return c.Status(fiber.StatusOK).JSON(response)
	}
}
