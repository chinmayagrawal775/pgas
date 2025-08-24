package main

import (
	"fmt"
	"pgas/pkg/processor"
	"pgas/pkg/providers"
	"pgas/pkg/providers/mastercard"
	"pgas/pkg/providers/visa"
)

func main() {

	// Initialize payment providers
	mastercardProvider := mastercard.GetNewMasterCardPaymentProvider()
	visaProvider := visa.GetNewVisaPaymentProvider()

	// Initialize the payment processor
	paymentProcessor := processor.NewPaymentProcessor([]providers.Provider{mastercardProvider, visaProvider})

	// Example payment request
	paymentRequests := providers.PaymentRequest{
		Mode:        "visa",
		Amount:      100.00,
		Currency:    "USD",
		CardNumber:  "4111111111111111",
		ExpiryMonth: "12",
		ExpiryYear:  "2025",
		CVV:         "123",
	}

	res, err := paymentProcessor.ProcessPayment(paymentRequests)
	if err != nil {
		fmt.Printf("payment failed: %v", err)
	}

	fmt.Printf("Payment Response: %+v\n", res)
}
