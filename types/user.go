package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserStore interface {
  CreateUser(context context.Context, user *UserCreate) (*User, error)
}

type UserCreate struct {
  Email string `json:"email" db:"email"`
  Name string `json:"name" db:"name"`
}

type User struct {
  Id uuid.UUID `json:"id" db:"id"`
  CreatedAt time.Time `json:"created_at" db:"created_at"`
  *UserCreate
}
