package provider

// EnvProvider is an interface for getting environment variables
type EnvProvider interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
}

// DefaultEnvProvider is a default implementation of EnvProvider
type DefaultEnvProvider struct {
	// 현재는 하드코딩된 값을 반환하지만,
	// 나중에 .env 파일이나 다른 소스에서 값을 가져오도록 수정할 수 있습니다.
}

// NewDefaultEnvProvider creates a new instance of DefaultEnvProvider
func NewDefaultEnvProvider() EnvProvider {
	return &DefaultEnvProvider{}
}

// GetString returns a string value for the given key
func (p *DefaultEnvProvider) GetString(key string) string {
	// 임시로 하드코딩된 값들
	defaults := map[string]string{
		"JWT_SECRET": "your-secret-key",
		"JWT_EXPIRY": "24h",
	}

	if value, exists := defaults[key]; exists {
		return value
	}
	return ""
}

// GetInt returns an integer value for the given key
func (p *DefaultEnvProvider) GetInt(key string) int {
	// 임시로 하드코딩된 값들
	defaults := map[string]int{
		"JWT_EXPIRY_HOURS": 24,
	}

	if value, exists := defaults[key]; exists {
		return value
	}
	return 0
}

// GetBool returns a boolean value for the given key
func (p *DefaultEnvProvider) GetBool(key string) bool {
	// 임시로 하드코딩된 값들
	defaults := map[string]bool{
		"JWT_ENABLED": true,
	}

	if value, exists := defaults[key]; exists {
		return value
	}
	return false
}
