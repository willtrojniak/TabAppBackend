package models

type substitutionGroupBase struct {
	Name string `json:"name" db:"name" validate:"required,min=1,max=64"`
}

type SubstitutionGroupUpdate struct {
	substitutionGroupBase
	SubstitutionItemIds []int `json:"substitution_item_ids" db:"substitution_item_ids" validate:"required,dive,gte=1"`
}

type SubstitutionGroupCreate struct {
	SubstitutionGroupUpdate
	ShopId int `json:"shop_id" db:"shop_id" validate:"required,gte=1"`
}

type SubstitutionGroup struct {
	substitutionGroupBase
	Substitutions []ItemOverview `json:"substitutions" db:"substitutions" validate:"required,dive"`
	Id            int            `json:"id" db:"id" validate:"required,gte=1"`
}
