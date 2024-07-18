package types

type categoryBase struct {
	Name string `json:"name" db:"name" validate:"required,min=1,max=64"`
}

type CategoryUpdate struct {
	categoryBase
	Index   *int  `json:"index" db:"index" validate:"required"`
	ItemIds []int `json:"item_ids" db:"item_ids" validate:"required,dive,gte=1"`
}

type CategoryCreate struct {
	ShopId int `json:"shop_id" db:"shop_id" validate:"required,gte=1"`
	CategoryUpdate
}

type CategoryOverview struct {
	categoryBase
	Id int `json:"id" db:"id" validate:"required,gte=1"`
}

type Category struct {
	Id int `json:"id" db:"id" validate:"required,gte=1"`
	CategoryCreate
}
