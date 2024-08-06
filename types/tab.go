package types

import (
	"log"
	"reflect"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type TabStatus int

const (
	TAB_STATUS_PENDING TabStatus = iota
	TAB_STATUS_CONFIRMED
	TAB_STATUS_CLOSED
)

func (s TabStatus) String() string {
	switch s {
	case TAB_STATUS_PENDING:
		return "pending"
	case TAB_STATUS_CONFIRMED:
		return "confirmed"
	case TAB_STATUS_CLOSED:
		return "closed"
	default:
		return "unknown"
	}
}

type OrderCreate struct {
	Id       int  `json:"id" db:"id" validate:"required,gte=1"`
	Quantity *int `json:"quantity" db:"quantity" validate:"required,gte=0"`
}

type ItemOrderCreate struct {
	OrderCreate
	Variants []OrderCreate `json:"variants" db:"variants" validate:"required,dive"`
}

type BillOrderCreate struct {
	Items []ItemOrderCreate `json:"items" db:"items" validate:"required,dive"`
}

type Bill struct {
	Id        int         `json:"id" db:"id" validate:"required,gte=1"`
	StartTime time.Time   `json:"start_time" db:"start_time" validate:"required"`
	IsPaid    bool        `json:"is_paid" db:"is_paid" validate:"required"`
	Items     []ItemOrder `json:"items" db:"items" validate:"required"`
}

/*
		{
		  Bills: [
		    {
		      id: 1,
		      start_time: July 2, 6 AM,
		      is_paid: false,
		      orders: [
		        {
		          id: 3
		          name: "Latte"
		          price: 4.50,
		          quantity: 22,
		          variants: [
		            {
	                id: 3
	                name: "Small"
	                price: 0.00,
	                quantity: 10,
		            },
		            {
	                id: 3
	                name: "Large"
	                price: 0.50,
	                quantity: 12,
		            }
		          ]
		        }
		      ]
		    }
		  ]
		}
*/
type TabBase struct {
	PaymentMethod       string  `json:"payment_method" db:"payment_method" validate:"required,oneof='in person' 'chartstring'"`
	Organization        string  `json:"organization" db:"organization" validate:"required,min=3,max=64"`
	DisplayName         string  `json:"display_name" db:"display_name" validate:"required,min=3,max=64"`
	StartDate           Date    `json:"start_date" db:"start_date" validate:"required,future"`
	EndDate             Date    `json:"end_date" db:"end_date" validate:"required"`
	DailyStartTime      Time    `json:"daily_start_time" db:"daily_start_time" validate:"required"`
	DailyEndTime        Time    `json:"daily_end_time" db:"daily_end_time" validate:"required"`
	ActiveDaysOfWk      int8    `json:"active_days_of_wk" db:"active_days_of_wk"`
	DollarLimitPerOrder float32 `json:"dollar_limit_per_order" db:"dollar_limit_per_order" validate:"gte=0"`
	VerificationMethod  string  `json:"verification_method" db:"verification_method" validate:"required,oneof='specify' 'voucher' 'email'"`
	PaymentDetails      string  `json:"payment_details" db:"payment_details"`
	BillingIntervalDays int     `json:"billing_interval_days" db:"billing_interval_days" validate:"gte=1,lte=365"`
}

type TabUpdate struct {
	TabBase
	VerificationList []string `json:"verification_list" db:"verification_list" validate:"required,dive,required,email"`
}

type TabCreate struct {
	TabUpdate
	ShopId  int    `json:"shop_id" db:"shop_id" validate:"required,gte=1"`
	OwnerId string `json:"owner_id" db:"owner_id" validate:"required"`
}

type TabOverview struct {
	TabCreate
	Id             int      `json:"id" db:"id" validate:"required,gte=1"`
	PendingUpdates *TabBase `json:"pending_updates" db:"pending_updates"`
	Status         string   `json:"status" db:"status"`
}

type Tab struct {
	TabOverview
	Bills []Bill `json:"bills" db:"bills" validate:"required,dive"`
}

func (t *TabOverview) IsActive() bool {
	today := DateOf(time.Now())
	return t.Status == TAB_STATUS_CONFIRMED.String() && !t.StartDate.After(today.Date) && !t.EndDate.Before(today.Date)
}

func TabUpdateStructLevelValidation(sl validator.StructLevel) {
	data := sl.Current().Interface().(TabUpdate)

	if data.DailyEndTime.Duration < data.DailyStartTime.Duration {
		field, _ := reflect.ValueOf(data).Type().FieldByName("DailyEndTime")
		tag, ok := field.Tag.Lookup("json")
		if !ok {
			tag = field.Name
		}
		sl.ReportError(data.DailyEndTime, tag, field.Name, "endafterstart", "")
	}

	if data.EndDate.Before(data.StartDate.Date) {
		field, _ := reflect.ValueOf(data).Type().FieldByName("EndDate")
		tag, ok := field.Tag.Lookup("json")
		if !ok {
			tag = field.Name
		}
		sl.ReportError(data.EndDate, tag, field.Name, "endafterstart", "")
	}

	chartstringPattern, err := regexp.Compile(`^([A-z0-9]{5})[ |-]?([A-z0-9]{5})(?:(?:-|\s)([A-z0-9]{5})|([A-z0-9]{5}))?$`)
	if err != nil {
		log.Fatal("Failed to compile chartstring expression")
	}
	if data.PaymentMethod == "chartstring" && !chartstringPattern.MatchString(data.PaymentDetails) {
		field, _ := reflect.ValueOf(data).Type().FieldByName("PaymentDetails")
		tag, ok := field.Tag.Lookup("json")
		if !ok {
			tag = field.Name
		}
		sl.ReportError(data.PaymentDetails, tag, field.Name, "charstringformat", "")
	}
}
