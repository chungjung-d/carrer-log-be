package auth

import (
	"career-log-be/models"
	"career-log-be/utils/jwt"
	"errors"

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
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input",
			})
		}

		// Validate input
		if err := validate.Struct(input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Validation failed",
				"details": err.Error(),
			})
		}

		// 사용자 찾기
		var user models.User
		result := db.Where("email = ?", input.Email).First(&user)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid credentials",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}

		// 비밀번호 확인
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}

		// JWT 토큰 생성
		token, err := jwtUtils.GenerateToken(user.ID, user.Email)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not generate token",
			})
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
