package types

import (
	"github.com/google/uuid"
)

type User struct {
	UUID          uuid.UUID
	Email         string
	Username      string
	CreditBalance uint64
}
