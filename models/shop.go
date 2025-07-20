package models

import "time"

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

type ShopOverview struct {
	Id uint `json:"id" db:"id" validate:"required,gte=1"`
	ShopCreate
}

type Shop struct {
	ShopOverview
	Locations []Location `json:"locations" db:"locations"`
	Users     []ShopUser `json:"users" db:"users"`
	ShopSlackData
}

type GetShopsQueryParams struct {
	Limit     int
	Offset    int
	IsMember  *bool
	UserId    *string
	IsPending *bool
}

type ShopUserCreate struct {
	Email string  `json:"email" db:"email" validate:"required,email,max=64"`
	Roles *uint32 `json:"roles" db:"roles" validate:"required"`
}

type ShopUser struct {
	User
	Roles       uint32    `json:"roles" db:"roles"`
	IsConfirmed bool      `json:"confirmed" db:"confirmed"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	IsOwner     bool      `json:"is_owner" db:"is_owner"`
}

type LocationUpdate struct {
	Name string `json:"name" db:"name" validate:"required,min=1,max=64"`
}

type LocationCreate struct {
	ShopId int `json:"shop_id" db:"shop_id" validate:"required,gte=1"`
	LocationUpdate
}

type Location struct {
	Id uint `json:"id" db:"id" validate:"required,gte=1"`
	LocationUpdate
}

func (shop *Shop) ConfirmedUsers() []*User {
	var users []*User
	for _, u := range shop.Users {
		if u.IsConfirmed {
			users = append(users, &u.User)
		}
	}
	return users
}

type ShopSlackData struct {
	SlackAccessToken *Token `json:"-" db:"slack_access_token"`
	SlackIntegrated  bool   `json:"slack_integrated" db:"slack_integrated"`
	ShopSlackDataUpdate
}

type ShopSlackDataUpdate struct {
	DailyUpdateSlackChannel    string `json:"daily_update_slack_channel" db:"daily_update_slack_channel" validate:"max=64"`         // Empty value indicates disabled
	TabRequestSlackChannel     string `json:"tab_request_slack_channel" db:"tab_request_slack_channel" validate:"max=64"`           // Empty value indicates disabled
	TabBillReceiptSlackChannel string `json:"tab_bill_receipt_slack_channel" db:"tab_bill_receipt_slack_channel" validate:"max=64"` // Empty value indicates disabled
}
