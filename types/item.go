package types

import "github.com/google/uuid"

type ItemOverview struct {
	Name      string    `json:"name" db:"name" validate:"required,min=1,max=64"`
	BasePrice *float32  `json:"base_price" db:"base_price" validate:"required,gte=0"`
	ShopId    uuid.UUID `json:"shop_id" db:"shop_id" validate:"required,uuid4"`
	Id        uuid.UUID `json:"id" db:"id" validate:"required,uuid4"`
}

type ItemUpdate struct {
	Name        string      `json:"name" db:"name" validate:"required,min=1,max=64"`
	BasePrice   *float32    `json:"base_price" db:"base_price" validate:"required,gte=0"`
	CategoryIds []uuid.UUID `json:"category_ids" db:"category_ids" validate:"required,dive,uuid4"`
}

type ItemCreate struct {
	ItemUpdate
	ShopId uuid.UUID `json:"shop_id" db:"shop_id" validate:"required,uuid4"`
}

type Item struct {
	ItemCreate
	Id uuid.UUID `json:"id" db:"id" validate:"required,uuid4"`
}

func (item *Item) GetOverview() ItemOverview {
	return ItemOverview{
		Id:        item.Id,
		ShopId:    item.Id,
		Name:      item.Name,
		BasePrice: item.BasePrice,
	}
}
