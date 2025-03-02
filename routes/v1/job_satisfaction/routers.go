package job_satisfaction

import (
	"career-log-be/middleware"
	job_satisfaction "career-log-be/services/job_satisfaction"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router) {
	jobSatisfactionRouter := router.Group("/job-satisfaction")

	// 보호된 라우트 그룹
	protected := jobSatisfactionRouter.Use(middleware.AuthMiddleware())

	// 직무 만족도 중요도 생성
	protected.Post("/importance", job_satisfaction.HandleCreateJobSatisfactionImportance())

	// 직무 만족도 초기화
	protected.Post("/init", job_satisfaction.HandleInitializeJobSatisfaction())

	// 현재 직무 만족도 조회
	protected.Get("/current", job_satisfaction.HandleGetCurrentJobSatisfaction())
}
