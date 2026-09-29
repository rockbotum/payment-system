package main

import (
	"fmt"

	paymentsystem "PaymentSystemModule"

	"github.com/shopspring/decimal"
)

func main() {
	processor := paymentsystem.NewPaymentProcessor()

	// Оплата кредитной картой
	card := paymentsystem.NewCreditCard("1234567890123456", "12/12/2027")
	fmt.Printf("Оплата через %s\n", card.GetName())

	err := processor.ProcessPayment(card, decimal.NewFromFloat(100.00))
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}

	// Оплата через PayPal
	paypal := paymentsystem.NewPayPal("user@example.com")
	fmt.Printf("\nОплата через %s\n", paypal.GetName())

	err = processor.ProcessPayment(paypal, decimal.NewFromFloat(50.00))
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}

	// Невалидная карта
	invalidCard := paymentsystem.NewCreditCard("1234", "12/2020")
	fmt.Printf("\nПопытка оплаты через %s с невалидными данными\n", invalidCard.GetName())

	err = processor.ProcessPayment(invalidCard, decimal.NewFromFloat(10.00))
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	}

	// История платежей
	fmt.Println("\n=== История платежей ===")
	for _, record := range processor.GetHistory() {
		fmt.Printf("Платеж %s: метод=%s, сумма=%v, комиссия=%v, итого=%v, успех=%v\n",
			record.ID, record.Method, record.Amount, record.Fee, record.Total, record.Success)
	}
}
