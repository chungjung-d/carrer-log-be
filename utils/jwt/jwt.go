package jwt

import (
	"career-log-be/config/env"
	"career-log-be/config/env/provider"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTUtils struct {
	config *env.JWTConfig
}

func NewJWTUtils(envProvider provider.EnvProvider) *JWTUtils {
	return &JWTUtils{
		config: env.NewJWTConfig(envProvider),
	}
}

type UserClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token for a user
func (j *JWTUtils) GenerateToken(userID, email string) (string, error) {
	claims := UserClaims{
		ID:    userID,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.config.GetExpiryDuration())),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(j.config.SecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return signedToken, nil
}

// ValidateToken validates the JWT token and returns the claims
func (j *JWTUtils) ValidateToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.config.SecretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
