package utils

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

func ValidateUsername(username string, maxLength int) error {
	if username == "" {
		return &ValidationError{Field: "username", Message: "Имя пользователя не может быть пустым"}
	}

	if utf8.RuneCountInString(username) > maxLength {
		return &ValidationError{Field: "username", Message: "Имя пользователя слишком длинное"}
	}

	invalidChars := regexp.MustCompile(`[<>'"&]`)
	if invalidChars.MatchString(username) {
		return &ValidationError{Field: "username", Message: "Имя пользователя содержит недопустимые символы"}
	}

	if strings.TrimSpace(username) == "" {
		return &ValidationError{Field: "username", Message: "Имя пользователя не может состоять только из пробелов"}
	}

	return nil
}

func ValidateRoomName(roomName string, maxLength int) error {
	if roomName == "" {
		return &ValidationError{Field: "roomName", Message: "Название комнаты не может быть пустым"}
	}

	if utf8.RuneCountInString(roomName) > maxLength {
		return &ValidationError{Field: "roomName", Message: "Название комнаты слишком длинное"}
	}

	invalidChars := regexp.MustCompile(`[<>'"&]`)
	if invalidChars.MatchString(roomName) {
		return &ValidationError{Field: "roomName", Message: "Название комнаты содержит недопустимые символы"}
	}

	if strings.TrimSpace(roomName) == "" {
		return &ValidationError{Field: "roomName", Message: "Название комнаты не может состоять только из пробелов"}
	}

	return nil
}

func SanitizeString(input string) string {
	htmlEscaped := strings.ReplaceAll(input, "<", "&lt;")
	htmlEscaped = strings.ReplaceAll(htmlEscaped, ">", "&gt;")
	htmlEscaped = strings.ReplaceAll(htmlEscaped, "&", "&amp;")
	htmlEscaped = strings.ReplaceAll(htmlEscaped, "\"", "&quot;")
	htmlEscaped = strings.ReplaceAll(htmlEscaped, "'", "&#39;")

	return htmlEscaped
}

func GenerateRoomID(roomName string) string {
	roomID := strings.ToLower(roomName)
	roomID = strings.ReplaceAll(roomID, " ", "-")

	reg := regexp.MustCompile(`[^a-z0-9\-]`)
	roomID = reg.ReplaceAllString(roomID, "")

	if len(roomID) > 50 {
		roomID = roomID[:50]
	}

	return roomID
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
