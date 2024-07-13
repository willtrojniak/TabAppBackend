package types

import "github.com/google/uuid"

type itemBase struct {
	Name      string   `json:"name" db:"name" validate:"required,min=1,max=64"`
	BasePrice *float32 `json:"base_price" db:"base_price" validate:"required,gte=0"`
}

type ItemUpdate struct {
	itemBase
	CategoryIds []uuid.UUID `json:"category_ids" db:"category_ids" validate:"required,dive,uuid4"`
	AddonIds    []uuid.UUID `json:"addon_ids" db:"addon_ids" validate:"required,dive,uuid4"`
}

type ItemCreate struct {
	ItemUpdate
	ShopId uuid.UUID `json:"shop_id" db:"shop_id" validate:"required,uuid4"`
}

type ItemOverview struct {
	itemBase
	Id uuid.UUID `json:"id" db:"id" validate:"required,uuid4"`
}

type Item struct {
	ItemOverview
	Categories []CategoryOverview `json:"categories" db:"categories" validate:"required,dive"`
	Variants   []ItemVariant      `json:"variants" db:"variants" validate:"required,dive"`
	Addons     []ItemOverview     `json:"addons" db:"addons" validate:"required,dive"`
}

func (item *Item) GetOverview() ItemOverview {
	return ItemOverview{
		Id: item.Id,
		itemBase: itemBase{
			Name:      item.Name,
			BasePrice: item.BasePrice,
		},
	}
}

type ItemVariantUpdate struct {
	Name  string   `json:"name" db:"name" validate:"required,min=1,max=64"`
	Price *float32 `json:"price" db:"price" validate:"required,gte=0"`
	Index *int     `json:"index" db:"index" validate:"required"`
}

type ItemVariantCreate struct {
	ItemVariantUpdate
	ShopId uuid.UUID `json:"shop_id" db:"shop_id" validate:"required,uuid4"`
	ItemId uuid.UUID `json:"item_id" db:"item_id" validate:"required,uuid4"`
}

type ItemVariant struct {
	ItemVariantUpdate
	Id uuid.UUID `json:"id" db:"id" validate:"required,uuid4"`
}
