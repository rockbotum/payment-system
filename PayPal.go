package paymentsystem

import (
	"fmt"
	"regexp"

	"github.com/shopspring/decimal"
)

type PayPal struct {
	email string
}

func NewPayPal(email string) *PayPal {
	return &PayPal{
		email: email,
	}
}

func (p *PayPal) GetName() string {
	return "PayPal"
}

func (p *PayPal) GetFree(amount decimal.Decimal) decimal.Decimal {
	return amount.Mul(decimal.NewFromFloat(0.03)).Add(decimal.NewFromFloat(0.5))
}

func (p *PayPal) Validate() bool {
	if len(p.email) > 0 && regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(p.email) {
		return true
	}
	return false
}

func (p *PayPal) Pay(amount decimal.Decimal) error {
	if !p.Validate() {
		return fmt.Errorf("Invalid PayPal email")
	} else {
		fmt.Printf("Processing PayPal payment of $%v\n", amount)
	}
	return nil
}
