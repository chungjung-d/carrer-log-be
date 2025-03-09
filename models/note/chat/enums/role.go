package enums

// MessageRole은 채팅 메시지의 역할을 나타내는 타입입니다
type MessageRole string

const (
	// UserRole은 사용자가 보낸 메시지를 나타냅니다
	UserRole MessageRole = "user"
	// AssistantRole은 AI 어시스턴트가 보낸 메시지를 나타냅니다
	AssistantRole MessageRole = "assistant"
)

// String은 MessageRole을 문자열로 변환합니다
func (r MessageRole) String() string {
	return string(r)
}

// IsValid는 MessageRole이 유효한 값인지 검사합니다
func (r MessageRole) IsValid() bool {
	switch r {
	case UserRole, AssistantRole:
		return true
	}
	return false
}
