package chat

import (
	"career-log-be/models/note/chat/enums"
	"career-log-be/utils"
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
