package env

import (
	"career-log-be/config/env/provider"
	"time"
)

type JWTConfig struct {
	SecretKey     []byte
	ExpiryHours   int
	SigningMethod string
}

func NewJWTConfig(envProvider provider.EnvProvider) *JWTConfig {
	return &JWTConfig{
		SecretKey:     []byte(envProvider.GetString("JWT_SECRET")),
		ExpiryHours:   envProvider.GetInt("JWT_EXPIRY_HOURS"),
		SigningMethod: "HS256",
	}
}

func (c *JWTConfig) GetExpiryDuration() time.Duration {
	return time.Duration(c.ExpiryHours) * time.Hour
}
