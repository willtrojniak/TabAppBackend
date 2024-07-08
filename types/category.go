package types

import "github.com/google/uuid"

type CategoryUpdate struct {
	Name  string `json:"name" db:"name" validate:"required,min=1,max=64"`
	Index int    `json:"index" db:"index" validate:"required"`
}

type CategoryCreate struct {
	ShopId uuid.UUID `json:"shop_id" db:"shop_id" validate:"required,uuid4"`
	CategoryUpdate
}

type Category struct {
	Id uuid.UUID `json:"id" db:"id" validate:"required,uuid4"`
	CategoryCreate
}
