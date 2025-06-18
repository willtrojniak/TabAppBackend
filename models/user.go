package models

import "time"

type UserUpdate struct {
	PreferredName *string `json:"preferred_name" db:"preferred_name" validate:"omitempty,min=2,max=64"`
}

type UserCreate struct {
	Id    string `json:"id" db:"id" validate:"required,max=255"`
	Email string `json:"email" db:"email" validate:"required,email,max=255"`
	Name  string `json:"name" db:"name" validate:"required,min=2,max=64"`
	UserUpdate
}

type User struct {
	UserCreate
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
