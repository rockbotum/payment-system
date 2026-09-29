package paymentsystem

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

var expiryDateRegexp = regexp.MustCompile(`^(0[1-9]|[12][0-9]|3[01])/(0[1-9]|1[0-2])/([0-9]{4})$`)

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
	matches := expiryDateRegexp.FindStringSubmatch(c.expiryDate)
	if len(c.number) != 16 || matches == nil {
		return false
	}

	c.day, c.month, c.year = matches[1], matches[2], matches[3]

	year, err := strconv.Atoi(c.year)
	if err != nil {
		return false
	}
	return year > time.Now().Year()
}

func (c *CreditCard) Pay(amount decimal.Decimal) error {
	if !c.Validate() {
		return fmt.Errorf("Invalid credit card details")
	} else {
		fmt.Printf("Processing credit card payment of $%v\n", amount)
	}
	return nil
}
