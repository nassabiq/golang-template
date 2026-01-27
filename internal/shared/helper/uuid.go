package helper

import "github.com/google/uuid"

type UUIDGeneratorInterface interface {
	GenerateID() string
}

type UUIDGenerator struct{}

func (u *UUIDGenerator) GenerateID() string {
	return uuid.NewString()
}
