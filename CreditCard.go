package paymentsystem

import (
	"fmt"
	"regexp"

	"github.com/shopspring/decimal"
)

type CreditCard struct {
	number     string
	expiryDate string
	day        string
	month      string
	year       string
}

func NewCreditCard(number string, expiryDate string) *CreditCard {
	return &CreditCard{
		number:     number,
		expiryDate: expiryDate,
	}
}

func (c *CreditCard) GetName() string {
	return "Credit Card"
}

func (c *CreditCard) GetFree(amount decimal.Decimal) decimal.Decimal {
	return amount.Mul(decimal.NewFromFloat(0.02))
}

func (c *CreditCard) Validate() bool {
	if len(c.number) == 16 && regexp.MustCompile(`^(0[1-9]|[12][0-9]|3[01])/(0[1-9]|1[0-2])/([0-9]{4})$`).MatchString(c.expiryDate) {
		matches := regexp.MustCompile(`^(0[1-9]|[12][0-9]|3[01])/(0[1-9]|1[0-2])/([0-9]{4})$`).FindStringSubmatch(c.expiryDate)
		c.day = matches[1]
		c.month = matches[2]
		c.year = matches[3]
		if c.year <= "2025" {
			return false
		} else {
			return true
		}
	}
	return false
}

func (c *CreditCard) Pay(amount decimal.Decimal) error {
	if !c.Validate() {
		return fmt.Errorf("Invalid credit card details")
	} else {
		fmt.Printf("Processing credit card payment of $%v\n", amount)
	}
	return nil
}
