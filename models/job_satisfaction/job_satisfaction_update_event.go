package job_satisfaction

import (
	"career-log-be/models/job_satisfaction/enums"
	"career-log-be/utils"

	"time"

	"gorm.io/gorm"
)

const (
	// JobSatisfactionEventPrefix는 직무 만족도 변경 이벤트 ID의 접두사입니다
	JobSatisfactionEventPrefix = "JOB_SAT_UPDATE_EVENT"
)

type JobSatisfactionUpdateEvent struct {
	ID                string                               `json:"id" gorm:"primaryKey;not null"`
	UserID            string                               `json:"userId" gorm:"index;not null"`
	EventType         enums.JobSatisfactionUpdateEventType `json:"eventType" gorm:"type:varchar(20);not null"`
	Workload          float64                              `json:"workload" gorm:"check:workload >= -100 AND workload <= 100;not null"`
	Compensation      float64                              `json:"compensation" gorm:"check:compensation >= -100 AND compensation <= 100;not null"`
	Growth            float64                              `json:"growth" gorm:"check:growth >= -100 AND growth <= 100;not null"`
	WorkEnvironment   float64                              `json:"workEnvironment" gorm:"check:work_environment >= -100 AND work_environment <= 100;column:work_environment;not null"`
	WorkRelationships float64                              `json:"workRelationships" gorm:"check:work_relationships >= -100 AND work_relationships <= 100;column:work_relationships;not null"`
	WorkValues        float64                              `json:"workValues" gorm:"check:work_values >= -100 AND work_values <= 100;column:work_values;not null"`
	SourceId          *string                              `json:"sourceId"` // 참조 ID (S3에 저장된 대화 내용 참조), nullable
	CreatedAt         time.Time                            `json:"createdAt" gorm:"not null"`
	UpdatedAt         time.Time                            `json:"updatedAt" gorm:"not null"`
}

func (u *JobSatisfactionUpdateEvent) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = utils.GenerateID(JobSatisfactionEventPrefix)
	return
}
