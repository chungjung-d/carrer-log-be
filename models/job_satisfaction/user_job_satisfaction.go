package job_satisfaction

import (
	"career-log-be/utils"
	"time"

	"gorm.io/gorm"
)

const (
	JobSatisfactionPrefix = "USR_JOB_SAT"
)

type UserJobSatisfaction struct {
	ID     string `json:"id" gorm:"primaryKey"`
	UserID string `json:"userId" gorm:"index"`

	// 만족도 점수
	Workload          float64 `json:"workload" gorm:"check:workload >= 0 AND workload <= 100"`
	Compensation      float64 `json:"compensation" gorm:"check:compensation >= 0 AND compensation <= 100"`
	Growth            float64 `json:"growth" gorm:"check:growth >= 0 AND growth <= 100"`
	WorkEnvironment   float64 `json:"workEnvironment" gorm:"check:work_environment >= 0 AND work_environment <= 100"`
	WorkRelationships float64 `json:"workRelationships" gorm:"check:work_relationships >= 0 AND work_relationships <= 100"`
	WorkValues        float64 `json:"workValues" gorm:"check:work_values >= 0 AND work_values <= 100"`

	// 중요도 점수
	WorkloadImportance          float64 `json:"workloadImportance" gorm:"check:workload_importance >= 0 AND workload_importance <= 100"`
	CompensationImportance      float64 `json:"compensationImportance" gorm:"check:compensation_importance >= 0 AND compensation_importance <= 100"`
	GrowthImportance            float64 `json:"growthImportance" gorm:"check:growth_importance >= 0 AND growth_importance <= 100"`
	WorkEnvironmentImportance   float64 `json:"workEnvironmentImportance" gorm:"check:work_environment_importance >= 0 AND work_environment_importance <= 100"`
	WorkRelationshipsImportance float64 `json:"workRelationshipsImportance" gorm:"check:work_relationships_importance >= 0 AND work_relationships_importance <= 100"`
	WorkValuesImportance        float64 `json:"workValuesImportance" gorm:"check:work_values_importance >= 0 AND work_values_importance <= 100"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (u *UserJobSatisfaction) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = utils.GenerateID(JobSatisfactionPrefix)
	return
}
