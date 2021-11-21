package v1

import (
	"github.com/google/uuid"
)

type UUIDService interface {
	NewUUID() string
}

type DefaultUUIDService struct{}

func (service DefaultUUIDService) NewUUID() string {
	return uuid.New().String()
}
