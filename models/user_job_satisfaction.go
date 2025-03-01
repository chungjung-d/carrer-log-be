package models

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
	Workload          int `json:"workload" gorm:"check:workload >= 0 AND workload <= 100"`
	Compensation      int `json:"compensation" gorm:"check:compensation >= 0 AND compensation <= 100"`
	Growth            int `json:"growth" gorm:"check:growth >= 0 AND growth <= 100"`
	WorkEnvironment   int `json:"workEnvironment" gorm:"check:work_environment >= 0 AND work_environment <= 100"`
	WorkRelationships int `json:"workRelationships" gorm:"check:work_relationships >= 0 AND work_relationships <= 100"`
	WorkValues        int `json:"workValues" gorm:"check:work_values >= 0 AND work_values <= 100"`

	// 중요도 점수
	WorkloadImportance          int `json:"workloadImportance" gorm:"check:workload_importance >= 0 AND workload_importance <= 100"`
	CompensationImportance      int `json:"compensationImportance" gorm:"check:compensation_importance >= 0 AND compensation_importance <= 100"`
	GrowthImportance            int `json:"growthImportance" gorm:"check:growth_importance >= 0 AND growth_importance <= 100"`
	WorkEnvironmentImportance   int `json:"workEnvironmentImportance" gorm:"check:work_environment_importance >= 0 AND work_environment_importance <= 100"`
	WorkRelationshipsImportance int `json:"workRelationshipsImportance" gorm:"check:work_relationships_importance >= 0 AND work_relationships_importance <= 100"`
	WorkValuesImportance        int `json:"workValuesImportance" gorm:"check:work_values_importance >= 0 AND work_values_importance <= 100"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (u *UserJobSatisfaction) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = utils.GenerateID(JobSatisfactionPrefix)
	return
}
