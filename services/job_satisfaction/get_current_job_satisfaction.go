package job_satisfaction

import (
	appErrors "career-log-be/errors"
	"career-log-be/models/job_satisfaction"
	utils "career-log-be/services/job_satisfaction/core/utils"
	"career-log-be/utils/response"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CurrentJobSatisfactionResponse struct {
	Workload          float64 `json:"workload"`
	Compensation      float64 `json:"compensation"`
	Growth            float64 `json:"growth"`
	WorkEnvironment   float64 `json:"workEnvironment"`
	WorkRelationships float64 `json:"workRelationships"`
	WorkValues        float64 `json:"workValues"`

	WorkloadImportance          float64 `json:"workloadImportance"`
	CompensationImportance      float64 `json:"compensationImportance"`
	GrowthImportance            float64 `json:"growthImportance"`
	WorkEnvironmentImportance   float64 `json:"workEnvironmentImportance"`
	WorkRelationshipsImportance float64 `json:"workRelationshipsImportance"`
	WorkValuesImportance        float64 `json:"workValuesImportance"`

	Score float64 `json:"score"`
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

		score := utils.CalculateWeightedScore(&satisfaction)

		resp := CurrentJobSatisfactionResponse{
			Workload:          satisfaction.Workload,
			Compensation:      satisfaction.Compensation,
			Growth:            satisfaction.Growth,
			WorkEnvironment:   satisfaction.WorkEnvironment,
			WorkRelationships: satisfaction.WorkRelationships,
			WorkValues:        satisfaction.WorkValues,

			WorkloadImportance:          satisfaction.WorkloadImportance,
			CompensationImportance:      satisfaction.CompensationImportance,
			GrowthImportance:            satisfaction.GrowthImportance,
			WorkEnvironmentImportance:   satisfaction.WorkEnvironmentImportance,
			WorkRelationshipsImportance: satisfaction.WorkRelationshipsImportance,
			WorkValuesImportance:        satisfaction.WorkValuesImportance,

			Score: score,
		}

		return response.Success(c, resp)
	}
}
