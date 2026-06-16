package paymentsystem

import (
	"fmt"
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

func (c *Crypto) GetFree(amount float64) float64 {
	return 0
}

func (c *Crypto) Validate() bool {
	if len(c.walletAddress) > 0 && fmt.Sprintf("%x", c.walletAddress) != "" {
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
