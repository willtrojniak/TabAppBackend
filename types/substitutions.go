package types

import "github.com/google/uuid"

type substitutionGroupBase struct {
	Name string `json:"name" db:"name" validate:"required,min=1,max=64"`
}

type SubstitutionGroupUpdate struct {
	substitutionGroupBase
	SubstitutionItemIds []uuid.UUID `json:"substitution_item_ids" db:"substitution_item_ids" validate:"required,dive,uuid4"`
}

type SubstitutionGroupCreate struct {
	SubstitutionGroupUpdate
	ShopId uuid.UUID `json:"shop_id" db:"shop_id" validate:"required,uuid4"`
}

type SubstitutionGroup struct {
	substitutionGroupBase
	Substitutions []ItemOverview `json:"substitutions" db:"substitutions" validate:"required,dive"`
	Id            uuid.UUID      `json:"id" db:"id" validate:"required,uuid4"`
}
