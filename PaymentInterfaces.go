package paymentsystem

import (
	"time"

	"github.com/shopspring/decimal"
)

type PaymentMethod interface {
	Pay(amount decimal.Decimal) error

	GetFree(amount decimal.Decimal) decimal.Decimal

	Validate() bool

	GetName() string
}

type PaymentProcessor interface {
	ProcessPayment(method PaymentMethod, amount decimal.Decimal) error

	GetHistory() []PaymentRecord
}

type PaymentRecord struct {
	ID        string
	Method    string
	Amount    decimal.Decimal
	Fee       decimal.Decimal
	Total     decimal.Decimal
	Success   bool
	Timestamp time.Time
}
