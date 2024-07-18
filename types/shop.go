package types

type PaymentMethod string

const (
	PaymentMethodInPerson    PaymentMethod = "in person"
	PaymentMethodChartstring PaymentMethod = "chartstring"
)

type ShopUpdate struct {
	Name           string   `json:"name" db:"name" validate:"required,min=1,max=64"`
	PaymentMethods []string `json:"payment_methods" db:"payment_methods" validate:"dive,oneof='in person' 'chartstring'"`
}

type ShopCreate struct {
	OwnerId string `json:"owner_id" db:"owner_id" validate:"required,max=255"`
	ShopUpdate
}

type Shop struct {
	Id uint `json:"id" db:"id" validate:"required,gte=1"`
	ShopCreate
}
