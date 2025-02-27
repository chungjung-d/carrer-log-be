package errors

import (
	"net/http"
)

// ErrorType은 에러의 종류를 나타냅니다
type ErrorType string

const (
	// 에러 타입 정의
	ErrorTypeValidation    ErrorType = "VALIDATION_ERROR"
	ErrorTypeAuthorization ErrorType = "AUTHORIZATION_ERROR"
	ErrorTypeNotFound      ErrorType = "NOT_FOUND_ERROR"
	ErrorTypeInternal      ErrorType = "INTERNAL_ERROR"
	ErrorTypeConflict      ErrorType = "CONFLICT_ERROR"
	ErrorTypeBadRequest    ErrorType = "BAD_REQUEST_ERROR"
)

// ErrorCode는 구체적인 에러 코드를 나타냅니다
type ErrorCode string

const (
	// 인증 관련 에러
	ErrorCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	ErrorCodeTokenRequired      ErrorCode = "TOKEN_REQUIRED"
	ErrorCodeInvalidToken       ErrorCode = "INVALID_TOKEN"
	ErrorCodeTokenExpired       ErrorCode = "TOKEN_EXPIRED"

	// 유효성 검사 관련 에러
	ErrorCodeInvalidInput  ErrorCode = "INVALID_INPUT"
	ErrorCodeRequiredField ErrorCode = "REQUIRED_FIELD"
	ErrorCodeInvalidFormat ErrorCode = "INVALID_FORMAT"

	// 리소스 관련 에러
	ErrorCodeResourceNotFound ErrorCode = "RESOURCE_NOT_FOUND"
	ErrorCodeResourceExists   ErrorCode = "RESOURCE_EXISTS"
	ErrorCodeResourceConflict ErrorCode = "RESOURCE_CONFLICT"

	// 서버 관련 에러
	ErrorCodeDatabaseError ErrorCode = "DATABASE_ERROR"
	ErrorCodeInternalError ErrorCode = "INTERNAL_SERVER_ERROR"
)

// 에러 타입에 따른 HTTP 상태 코드 매핑
var errorTypeToStatusCode = map[ErrorType]int{
	ErrorTypeValidation:    http.StatusBadRequest,
	ErrorTypeAuthorization: http.StatusUnauthorized,
	ErrorTypeNotFound:      http.StatusNotFound,
	ErrorTypeInternal:      http.StatusInternalServerError,
	ErrorTypeConflict:      http.StatusConflict,
	ErrorTypeBadRequest:    http.StatusBadRequest,
}
