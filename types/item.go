package types

type itemBase struct {
	Name      string   `json:"name" db:"name" validate:"required,min=1,max=64"`
	BasePrice *float32 `json:"base_price" db:"base_price" validate:"required,gte=0"`
}

type ItemUpdate struct {
	itemBase
	CategoryIds          []int `json:"category_ids" db:"category_ids" validate:"required,dive,gte=1"`
	AddonIds             []int `json:"addon_ids" db:"addon_ids" validate:"required,dive,gte=1"`
	SubstitutionGroupIds []int `json:"substitution_group_ids" db:"substitution_group_ids" validate:"required,dive,gte=1"`
}

type ItemCreate struct {
	ItemUpdate
	ShopId int `json:"shop_id" db:"shop_id" validate:"required,gte=1"`
}

type ItemOverview struct {
	itemBase
	Id int `json:"id" db:"id" validate:"required,gte=1"`
}

type ItemOrder struct {
	ItemOverview
	Variants []ItemVariant `json:"variants" db:"variants" validate:"required,dive"`
}

type Item struct {
	ItemOverview
	Categories         []CategoryOverview  `json:"categories" db:"categories" validate:"required,dive"`
	Variants           []ItemVariant       `json:"variants" db:"variants" validate:"required,dive"`
	Addons             []ItemOverview      `json:"addons" db:"addons" validate:"required,dive"`
	SubstitutionGroups []SubstitutionGroup `json:"substitution_groups" db:"substitution_groups" validate:"required,dive"`
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

type itemVariantBase struct {
	Name  string   `json:"name" db:"name" validate:"required,min=1,max=64"`
	Price *float32 `json:"price" db:"price" validate:"required,gte=0"`
}

type ItemVariantUpdate struct {
	itemVariantBase
	Index *int `json:"index" db:"index" validate:"required"`
}

type ItemVariantCreate struct {
	ItemVariantUpdate
	ShopId int `json:"shop_id" db:"shop_id" validate:"required,gte=1"`
	ItemId int `json:"item_id" db:"item_id" validate:"required,gte=1"`
}

type ItemVariant struct {
	itemVariantBase
	Id int `json:"id" db:"id" validate:"required,gte=1"`
}
