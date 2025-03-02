package job_satisfaction

import (
	appErrors "career-log-be/errors"
	"career-log-be/models/job_satisfaction"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CurrentJobSatisfactionResponse struct {
	Workload          int `json:"workload"`
	Compensation      int `json:"compensation"`
	Growth            int `json:"growth"`
	WorkEnvironment   int `json:"workEnvironment"`
	WorkRelationships int `json:"workRelationships"`
	WorkValues        int `json:"workValues"`
}

// HandleGetCurrentJobSatisfaction는 현재 사용자의 직무 만족도를 조회하는 핸들러입니다
func HandleGetCurrentJobSatisfaction() fiber.Handler {
	return func(c *fiber.Ctx) error {
		db := c.Locals("db").(*gorm.DB)
		userID := c.Locals("userID").(string)

		var satisfaction job_satisfaction.UserJobSatisfaction
		result := db.Where("user_id = ?", userID).First(&satisfaction)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				return appErrors.NewNotFoundError(
					appErrors.ErrorCodeResourceNotFound,
					"Not found job satisfaction data",
				)
			}
			return appErrors.NewInternalError(
				appErrors.ErrorCodeDatabaseError,
				"Error occurred while retrieving job satisfaction data",
				result.Error,
			)
		}

		resp := CurrentJobSatisfactionResponse{
			Workload:          satisfaction.Workload,
			Compensation:      satisfaction.Compensation,
			Growth:            satisfaction.Growth,
			WorkEnvironment:   satisfaction.WorkEnvironment,
			WorkRelationships: satisfaction.WorkRelationships,
			WorkValues:        satisfaction.WorkValues,
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    resp,
		})
	}
}
