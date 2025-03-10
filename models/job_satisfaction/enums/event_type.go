package enums

import "database/sql/driver"

type JobSatisfactionUpdateEventType string

const (
	InitEvent         JobSatisfactionUpdateEventType = "INIT_EVENT"
	ChatAnalysisEvent JobSatisfactionUpdateEventType = "CHAT_ANALYSIS_EVENT"
)

// Value - SQL을 위한 직렬화
func (et JobSatisfactionUpdateEventType) Value() (driver.Value, error) {
	return string(et), nil
}

// Scan - SQL에서 역직렬화
func (et *JobSatisfactionUpdateEventType) Scan(value interface{}) error {
	*et = JobSatisfactionUpdateEventType(value.(string))
	return nil
}

// IsValid - 이벤트 타입 유효성 검사
func (et JobSatisfactionUpdateEventType) IsValid() bool {
	switch et {
	case InitEvent, ChatAnalysisEvent:
		return true
	}
	return false
}

// String - 문자열 변환
func (et JobSatisfactionUpdateEventType) String() string {
	return string(et)
}
