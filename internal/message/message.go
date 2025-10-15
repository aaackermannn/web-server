package message

import (
	"html"
	"regexp"
	"time"
)

type Message struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	RoomID    string    `json:"roomId"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
}

type MessageRequest struct {
	Content string `json:"content"`
	RoomID  string `json:"roomId"`
}

type MessageResponse struct {
	Message *Message `json:"message,omitempty"`
	Error   string   `json:"error,omitempty"`
	Success bool     `json:"success"`
}

func NewMessage(id, userID, username, content, roomID, msgType string) *Message {
	return &Message{
		ID:        id,
		UserID:    userID,
		Username:  username,
		Content:   SanitizeContent(content),
		RoomID:    roomID,
		Timestamp: time.Now(),
		Type:      msgType,
	}
}

func SanitizeContent(content string) string {
	sanitized := html.EscapeString(content)
	re := regexp.MustCompile(`[<>'"&]`)
	sanitized = re.ReplaceAllString(sanitized, "")

	return sanitized
}

func ValidateMessage(content string, maxLength int) error {
	if len(content) == 0 {
		return &ValidationError{Field: "content", Message: "Сообщение не может быть пустым"}
	}

	if len(content) > maxLength {
		return &ValidationError{Field: "content", Message: "Сообщение слишком длинное"}
	}

	if regexp.MustCompile(`^\s*$`).MatchString(content) {
		return &ValidationError{Field: "content", Message: "Сообщение не может состоять только из пробелов"}
	}

	return nil
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
