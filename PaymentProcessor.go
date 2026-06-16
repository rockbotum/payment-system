package paymentsystem

import (
	"time"

	"github.com/shopspring/decimal"
)

type DefaultPaymentProcessor struct {
	history []PaymentRecord
}

func NewPaymentProcessor() PaymentProcessor {
	return &DefaultPaymentProcessor{
		history: make([]PaymentRecord, 0),
	}
}

func (p *DefaultPaymentProcessor) ProcessPayment(method PaymentMethod, amount decimal.Decimal) error {
	fee := method.GetFree(amount)
	total := decimal.Sum(amount, fee)
	err := method.Pay(amount)
	success := true
	if !method.Validate() {
		success = false
		return err
	}
	record := PaymentRecord{
		ID:        (string)(len(p.history) + 1),
		Method:    method.GetName(),
		Amount:    amount,
		Fee:       fee,
		Total:     total,
		Success:   success,
		Timestamp: time.Now(),
	}
	p.history = append(p.history, record)
	return nil
}

func (p *DefaultPaymentProcessor) GetHistory() []PaymentRecord {
	if len(p.history) == 0 {
		return nil
	}
	return p.history
}
