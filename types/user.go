package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserStore interface {
  CreateUser(context context.Context, user *UserCreate) (*uuid.UUID, error)
}

type UserCreate struct {
  Email string `json:"email" db:"email"`
  Name string `json:"name" db:"name"`
  PreferredName *string `json:"preferred_name" db:"preferred_name"`
}

type User struct {
  UserCreate
  Id uuid.UUID `json:"id" db:"id"`
  CreatedAt time.Time `json:"created_at" db:"created_at"`
}
