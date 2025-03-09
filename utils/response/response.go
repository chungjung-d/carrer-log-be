package response

import "github.com/gofiber/fiber/v2"

// Success는 성공 응답을 생성합니다
func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// Created는 리소스 생성 성공 응답을 생성합니다
func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// Accepted는 요청이 수락되었음을 나타내는 응답을 생성합니다
func Accepted(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": true,
		"message": message,
	})
}

// NoContent는 성공했지만 반환할 데이터가 없는 응답을 생성합니다
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}
