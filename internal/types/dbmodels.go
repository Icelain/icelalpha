package types

import (
	"github.com/google/uuid"
)

// Icealpha user schema
type User struct {

	// use uuid as primary key and main identity
	UUID     uuid.UUID
	Email    string
	Username string

	// keep track of credits expended by the user
	CreditBalance uint64
}
