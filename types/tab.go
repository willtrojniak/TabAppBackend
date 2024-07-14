package types

import (
	"time"

	"github.com/google/uuid"
)

type TabUpdate struct {
	PaymentMethod       string    `json:"payment_method" db:"payment_method" validate:"required,oneof='in person' 'chartstring'"`
	Organization        string    `json:"organization" db:"organization" validate:"required,min=3,max=64"`
	DisplayName         string    `json:"display_name" db:"display_name" validate:"required,min=3,max=64"`
	StartDate           time.Time `json:"start_date" db:"start_date" validate:"required"`
	EndDate             time.Time `json:"end_date" db:"end_date" validate:"required"`
	DailyStartTime      Time      `json:"daily_start_time" db:"daily_start_time" validate:"required"`
	DailyEndTime        Time      `json:"daily_end_time" db:"daily_end_time" validate:"required"`
	ActiveDaysOfWk      uint8     `json:"active_days_of_wk" db:"active_days_of_wk"`
	DollarLimitPerOrder float32   `json:"dollar_limit_per_order" db:"dollar_limit_per_order" validate:"gte=0"`
	VerificationMethod  string    `json:"verification_method" db:"verification_method" validate:"required,oneof='specify' 'voucher' 'email'"`
	PaymentDetails      string    `json:"payment_details" db:"payment_details"`
	BillingIntervalDays int       `json:"billing_interval_days" db:"billing_interval_days" validate:"gte=1,lte=365"`
	VerificationList    []string  `json:"verification_list" db:"verification_list" validate:"required,dive,required,email"`
}

type TabCreate struct {
	TabUpdate
	ShopId  uuid.UUID `json:"shop_id" db:"shop_id" validate:"required,uuid4"`
	OwnerId string    `json:"owner_id" db:"owner_id" validate:"required"`
}

type Tab struct {
	TabCreate
	Id     uint   `json:"id" db:"id" validate:"required,gte=1"`
	Status string `json:"status" db:"status"`
}
