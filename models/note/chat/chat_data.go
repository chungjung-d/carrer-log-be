package chat

import (
	"career-log-be/models/note/chat/enums"
	"career-log-be/utils"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

const (
	MessagePrefix = "MSG"
)

type Message struct {
	ID        string            `json:"id"`
	Role      enums.MessageRole `json:"role"`
	Content   string            `json:"content"`
	Timestamp time.Time         `json:"timestamp"`
}

type ChatMetadata struct {
	MessageCount  int       `json:"message_count"`
	LastMessageAt time.Time `json:"last_message_at"`
}

type ChatData struct {
	Messages []Message    `json:"messages"`
	Metadata ChatMetadata `json:"metadata"`
}

// NewChatData creates a new empty ChatData with initialized metadata
func NewChatData() ChatData {
	now := time.Now()
	return ChatData{
		Messages: []Message{},
		Metadata: ChatMetadata{
			MessageCount:  0,
			LastMessageAt: now,
		},
	}
}

// NewMessage creates a new message with generated ID
func NewMessage(role enums.MessageRole, content string) Message {
	return Message{
		ID:        utils.GenerateID(MessagePrefix),
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	}
}

// AddMessage adds a new message to ChatData and updates metadata
func (cd *ChatData) AddMessage(role enums.MessageRole, content string) {
	message := NewMessage(role, content)
	cd.Messages = append(cd.Messages, message)
	cd.Metadata.MessageCount++
	cd.Metadata.LastMessageAt = message.Timestamp
}

// Scan implements the sql.Scanner interface
func (cd *ChatData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", value)
	}

	return json.Unmarshal(bytes, cd)
}

// Value implements the driver.Valuer interface
func (cd ChatData) Value() (driver.Value, error) {
	return json.Marshal(cd)
}
