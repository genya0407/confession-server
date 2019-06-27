package utils

import (
	"github.com/google/uuid"
)

func MustNewUUID() uuid.UUID {
	u, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return u
}
