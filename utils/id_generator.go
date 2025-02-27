package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GenerateID creates a custom ID with format: PREFIX_TIMESTAMP_RANDOMSEED
func GenerateID(prefix string) string {
	// Get current timestamp in milliseconds
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	// Generate UUID and get first 6 characters (excluding hyphens)
	uuid := strings.ReplaceAll(uuid.New().String(), "-", "")
	randomSeed := uuid[:6]

	// Combine all parts
	return fmt.Sprintf("%s_%d_%s", prefix, timestamp, randomSeed)
}
