package types

import "github.com/google/uuid"

type ItemUpdate struct {
	Name      string   `json:"name" db:"name" validate:"required,min=1,max=64"`
	BasePrice *float32 `json:"base_price" db:"base_price" validate:"required,gte=0"`
}

type ItemCreate struct {
	ShopId uuid.UUID `json:"shop_id" db:"shop_id" validate:"required,uuid4"`
	ItemUpdate
}

type Item struct {
	Id uuid.UUID `json:"id" db:"id" validate:"required,uuid4"`
	ItemCreate
}
