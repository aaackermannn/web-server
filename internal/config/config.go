package config

import (
	"os"
)

type Config struct {
	Port              string
	MaxMessageLength  int
	MaxUsernameLength int
	MaxRoomNameLength int
	PingInterval      int
	PongWait          int
	WriteWait         int
}

func New() *Config {
	return &Config{
		Port:              getEnv("PORT", "8080"),
		MaxMessageLength:  500,
		MaxUsernameLength: 20,
		MaxRoomNameLength: 30,
		PingInterval:      54,
		PongWait:          60,
		WriteWait:         10,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
