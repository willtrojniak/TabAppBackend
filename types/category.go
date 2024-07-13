package types

import "github.com/google/uuid"

type categoryBase struct {
	Name string `json:"name" db:"name" validate:"required,min=1,max=64"`
}

type CategoryUpdate struct {
	categoryBase
	Index   *int        `json:"index" db:"index" validate:"required"`
	ItemIds []uuid.UUID `json:"item_ids" db:"item_ids" validate:"required,dive,uuid4"`
}

type CategoryCreate struct {
	ShopId uuid.UUID `json:"shop_id" db:"shop_id" validate:"required,uuid4"`
	CategoryUpdate
}

type CategoryOverview struct {
	categoryBase
	Id uuid.UUID `json:"id" db:"id" validate:"required,uuid4"`
}

type Category struct {
	Id uuid.UUID `json:"id" db:"id" validate:"required,uuid4"`
	CategoryCreate
}
