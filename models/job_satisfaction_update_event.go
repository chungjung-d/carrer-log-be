package models

import (
	"career-log-be/utils"
	"time"

	"gorm.io/gorm"
)

const (
	// JobSatisfactionEventPrefix는 직무 만족도 변경 이벤트 ID의 접두사입니다
	JobSatisfactionEventPrefix = "JOB_SAT_UPDATE_EVENT"
)

type JobSatisfactionUpdateEvent struct {
	ID                string    `json:"id" gorm:"primaryKey"`
	UserID            string    `json:"userId" gorm:"index"`
	Workload          int       `json:"workload"` // 변화량 (양수 또는 음수)
	Compensation      int       `json:"compensation"`
	Growth            int       `json:"growth"`
	WorkEnvironment   int       `json:"workEnvironment" gorm:"column:work_environment"`
	WorkRelationships int       `json:"workRelationships" gorm:"column:work_relationships"`
	WorkValues        int       `json:"workValues" gorm:"column:work_values"`
	SourceId          string    `json:"sourceId"` // 참조 ID (S3에 저장된 대화 내용 참조)
	CreatedAt         time.Time `json:"createdAt"`
}

func (u *JobSatisfactionUpdateEvent) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = utils.GenerateID(JobSatisfactionEventPrefix)
	return
}
