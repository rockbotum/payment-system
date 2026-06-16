package paymentsystem

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type Crypto struct {
	walletAddress string
}

func NewCrypto(walletAddress string) *Crypto {
	return &Crypto{
		walletAddress: walletAddress,
	}
}

func (c *Crypto) GetName() string {
	return "Cryptocurrency"
}

func (c *Crypto) GetFree(amount decimal.Decimal) decimal.Decimal {
	return amount.Mul(decimal.NewFromFloat(0.5))
}

func (c *Crypto) Validate() bool {

	if strings.HasPrefix(c.walletAddress, "0x") && len(c.walletAddress) == 42 {
		return true
	}
	return false
}

func (c *Crypto) Pay(amount float64) error {
	if !c.Validate() {
		return fmt.Errorf("Invalid cryptocurrency wallet address")
	} else {
		fmt.Printf("Processing cryptocurrency payment of $%.2f\n", amount)
	}
	return nil
}
