package util

import "github.com/google/uuid"

func GenerateTaskID() string {
	return uuid.New().String()
}
