package models

import (
	"time"
)

type UserJobSatisfactionImportance struct {
	ID                string    `json:"id" gorm:"primaryKey"`
	Workload          int       `json:"workload" gorm:"check:workload >= 0 AND workload <= 100"`
	Compensation      int       `json:"compensation" gorm:"check:compensation >= 0 AND compensation <= 100"`
	Growth            int       `json:"growth" gorm:"check:growth >= 0 AND growth <= 100"`
	WorkEnvironment   int       `json:"workEnvironment" gorm:"check:work_environment >= 0 AND work_environment <= 100"`
	WorkRelationships int       `json:"workRelationships" gorm:"check:work_relationships >= 0 AND work_relationships <= 100"`
	WorkValues        int       `json:"workValues" gorm:"check:work_values >= 0 AND work_values <= 100"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
