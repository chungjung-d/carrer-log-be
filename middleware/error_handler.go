package middleware

import (
	appErrors "career-log-be/errors"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler는 애플리케이션 에러를 HTTP 응답으로 변환합니다
func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// 기본 에러 응답
		code := fiber.StatusInternalServerError
		message := "Internal Server Error"
		errorCode := appErrors.ErrorCodeInternalError
		var details any

		// AppError 타입인지 확인
		var appError *appErrors.AppError
		if errors.As(err, &appError) {
			code = appError.StatusCode()
			message = appError.Message
			errorCode = appError.Code
			details = appError.Details

			// 내부 에러는 로깅
			if appError.Type == appErrors.ErrorTypeInternal {
				log.Printf("Internal error: %s\nDebug info: %s", appError.Error(), appError.DebugInfo)
			}
		} else {
			// 일반 에러는 로깅하고 일반적인 메시지 반환
			log.Printf("Unhandled error: %v", err)
		}

		// JSON 응답 반환
		return c.Status(code).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    errorCode,
				"message": message,
				"details": details,
			},
		})
	}
}
