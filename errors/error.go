package errors

import (
	"fmt"
)

// AppError는 애플리케이션에서 사용하는 에러 구조체입니다
type AppError struct {
	Type      ErrorType `json:"-"`                 // 에러 타입 (내부용)
	Code      ErrorCode `json:"code"`              // 에러 코드 (클라이언트용)
	Message   string    `json:"message"`           // 에러 메시지
	Details   any       `json:"details,omitempty"` // 추가 정보
	DebugInfo string    `json:"-"`                 // 디버깅 정보 (로그용)
	Err       error     `json:"-"`                 // 원본 에러
}

// Error는 error 인터페이스를 구현합니다
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap은 원본 에러를 반환합니다
func (e *AppError) Unwrap() error {
	return e.Err
}

// StatusCode는 에러 타입에 따른 HTTP 상태 코드를 반환합니다
func (e *AppError) StatusCode() int {
	if code, exists := errorTypeToStatusCode[e.Type]; exists {
		return code
	}
	return 500 // 기본값은 500 Internal Server Error
}

// NewValidationError는 유효성 검사 에러를 생성합니다
func NewValidationError(code ErrorCode, message string, details any) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Code:    code,
		Message: message,
		Details: details,
	}
}

// NewAuthorizationError는 인증 관련 에러를 생성합니다
func NewAuthorizationError(code ErrorCode, message string) *AppError {
	return &AppError{
		Type:    ErrorTypeAuthorization,
		Code:    code,
		Message: message,
	}
}

// NewNotFoundError는 리소스를 찾을 수 없는 에러를 생성합니다
func NewNotFoundError(code ErrorCode, message string) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Code:    code,
		Message: message,
	}
}

// NewConflictError는 리소스 충돌 에러를 생성합니다
func NewConflictError(code ErrorCode, message string) *AppError {
	return &AppError{
		Type:    ErrorTypeConflict,
		Code:    code,
		Message: message,
	}
}

// NewInternalError는 내부 서버 에러를 생성합니다
func NewInternalError(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Type:      ErrorTypeInternal,
		Code:      code,
		Message:   message,
		Err:       err,
		DebugInfo: fmt.Sprintf("%+v", err),
	}
}

// NewBadRequestError는 잘못된 요청 에러를 생성합니다
func NewBadRequestError(code ErrorCode, message string) *AppError {
	return &AppError{
		Type:    ErrorTypeBadRequest,
		Code:    code,
		Message: message,
	}
}
