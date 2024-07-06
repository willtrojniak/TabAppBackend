package types

import "github.com/google/uuid"

type PaymentMethod string

const (
	PaymentMethodInPerson    PaymentMethod = "in person"
	PaymentMethodChartstring PaymentMethod = "chartstring"
)

type ShopCreate struct {
	OwnerId string `json:"owner_id" db:"owner_id" validate:"required,max=255"`
	Name    string `json:"name" db:"name" validate:"required,min=1,max=64"`
}

type ShopUpdate struct {
	Id uuid.UUID `json:"id" db:"id" validate:"required"`
	ShopCreate
}

type Shop struct {
	ShopUpdate
}
