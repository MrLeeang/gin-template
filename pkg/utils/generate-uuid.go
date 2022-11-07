package utils

import (
	uuid "github.com/satori/go.uuid"
)

func GetUuid() string {
	return uuid.NewV4().String()
}
