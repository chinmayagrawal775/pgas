package processor

import (
	"testing"

	"pgas/pkg/providers"
	"pgas/pkg/providers/mastercard"
	"pgas/pkg/providers/visa"
)

func TestNewPaymentProcessor(t *testing.T) {
	// Test with empty providers
	processor := NewPaymentProcessor([]providers.Provider{})
	if processor == nil {
		t.Fatal("Expected processor to be created")
	}

	// Test with providers
	mastercardProvider := mastercard.GetNewMasterCardPaymentProvider()
	visaProvider := visa.GetNewVisaPaymentProvider()

	processor = NewPaymentProcessor([]providers.Provider{mastercardProvider, visaProvider})
	if processor == nil {
		t.Fatal("Expected processor to be created with providers")
	}
}

func TestGetProvider(t *testing.T) {
	mastercardProvider := mastercard.GetNewMasterCardPaymentProvider()
	visaProvider := visa.GetNewVisaPaymentProvider()

	processor := NewPaymentProcessor([]providers.Provider{mastercardProvider, visaProvider})

	// Test valid provider
	provider, err := processor.getProvider("mastercard")
	if err != nil {
		t.Fatalf("Expected no error for valid provider, got: %v", err)
	}
	if provider == nil {
		t.Fatal("Expected provider to be returned")
	}
	if provider.GetName() != "mastercard" {
		t.Errorf("Expected provider name 'mastercard', got: %s", provider.GetName())
	}

	// Test invalid provider
	_, err = processor.getProvider("invalid_provider")
	if err == nil {
		t.Fatal("Expected error for invalid provider")
	}
	expectedError := "invalid provider name provided: 'invalid_provider'"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got: '%s'", expectedError, err.Error())
	}
}

func TestProcessPayment_Success(t *testing.T) {
	mastercardProvider := mastercard.GetNewMasterCardPaymentProvider()
	visaProvider := visa.GetNewVisaPaymentProvider()

	processor := NewPaymentProcessor([]providers.Provider{mastercardProvider, visaProvider})

	// Test successful payment with Visa
	request := providers.PaymentRequest{
		Mode:        "visa",
		Amount:      100.00,
		Currency:    "USD",
		CardNumber:  "4111111111111111",
		ExpiryMonth: "12",
		ExpiryYear:  "2025",
		CVV:         "123",
	}

	response, err := processor.ProcessPayment(request)
	if err != nil {
		t.Fatalf("Expected successful payment, got error: %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}

	if response.TransactionID == "" {
		t.Error("Expected transaction ID to be set")
	}

	if response.Status == "" {
		t.Error("Expected status to be set")
	}

	if response.Amount != request.Amount {
		t.Errorf("Expected amount %f, got %f", request.Amount, response.Amount)
	}

	if response.Currency != request.Currency {
		t.Errorf("Expected currency %s, got %s", request.Currency, response.Currency)
	}
}

func TestProcessPayment_InvalidProvider(t *testing.T) {
	mastercardProvider := mastercard.GetNewMasterCardPaymentProvider()
	processor := NewPaymentProcessor([]providers.Provider{mastercardProvider})

	request := providers.PaymentRequest{
		Mode:        "invalid_provider",
		Amount:      100.00,
		Currency:    "USD",
		CardNumber:  "4111111111111111",
		ExpiryMonth: "12",
		ExpiryYear:  "2025",
		CVV:         "123",
	}

	_, err := processor.ProcessPayment(request)
	if err == nil {
		t.Fatal("Expected error for invalid provider")
	}

	if err.ErrorCode != "INVALID_PROVIDER" {
		t.Errorf("Expected error code 'INVALID_PROVIDER', got: %s", err.ErrorCode)
	}
}

func TestProcessPayment_ValidationError(t *testing.T) {
	mastercardProvider := mastercard.GetNewMasterCardPaymentProvider()
	processor := NewPaymentProcessor([]providers.Provider{mastercardProvider})

	// Test with invalid amount
	request := providers.PaymentRequest{
		Mode:        "mastercard",
		Amount:      0, // Invalid amount
		Currency:    "USD",
		CardNumber:  "5555555555554444",
		ExpiryMonth: "12",
		ExpiryYear:  "2025",
		CVV:         "123",
	}

	_, err := processor.ProcessPayment(request)
	if err == nil {
		t.Fatal("Expected error for invalid amount")
	}

	if err.ErrorCode != "INVALID_REQUEST" {
		t.Errorf("Expected error code 'INVALID_REQUEST', got: %s", err.ErrorCode)
	}
}

func TestProcessPayment_EmptyCardNumber(t *testing.T) {
	mastercardProvider := mastercard.GetNewMasterCardPaymentProvider()
	processor := NewPaymentProcessor([]providers.Provider{mastercardProvider})

	request := providers.PaymentRequest{
		Mode:        "mastercard",
		Amount:      100.00,
		Currency:    "USD",
		CardNumber:  "", // Empty card number
		ExpiryMonth: "12",
		ExpiryYear:  "2025",
		CVV:         "123",
	}

	_, err := processor.ProcessPayment(request)
	if err == nil {
		t.Fatal("Expected error for empty card number")
	}

	if err.ErrorCode != "INVALID_REQUEST" {
		t.Errorf("Expected error code 'INVALID_REQUEST', got: %s", err.ErrorCode)
	}
}

func TestProcessPayment_InvalidCVV(t *testing.T) {
	mastercardProvider := mastercard.GetNewMasterCardPaymentProvider()
	processor := NewPaymentProcessor([]providers.Provider{mastercardProvider})

	request := providers.PaymentRequest{
		Mode:        "mastercard",
		Amount:      100.00,
		Currency:    "USD",
		CardNumber:  "5555555555554444",
		ExpiryMonth: "12",
		ExpiryYear:  "2025",
		CVV:         "12", // Invalid CVV (too short)
	}

	_, err := processor.ProcessPayment(request)
	if err == nil {
		t.Fatal("Expected error for invalid CVV")
	}

	if err.ErrorCode != "INVALID_REQUEST" {
		t.Errorf("Expected error code 'INVALID_REQUEST', got: %s", err.ErrorCode)
	}
}

func TestProcessPayment_EdgeCaseAmounts(t *testing.T) {
	mastercardProvider := mastercard.GetNewMasterCardPaymentProvider()
	processor := NewPaymentProcessor([]providers.Provider{mastercardProvider})

	testCases := []struct {
		name   string
		amount float64
		valid  bool
	}{
		{"minimum amount", 0.01, true},
		{"small amount", 1.00, true},
		{"large amount", 999999.99, true},
		{"zero amount", 0.00, false},
		{"negative amount", -1.00, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := providers.PaymentRequest{
				Mode:        "mastercard",
				Amount:      tc.amount,
				Currency:    "USD",
				CardNumber:  "5555555555554444",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "123",
			}

			_, err := processor.ProcessPayment(request)

			if tc.valid && err != nil {
				t.Errorf("Expected success for amount %f, got error: %v", tc.amount, err)
			}

			if !tc.valid && err == nil {
				t.Errorf("Expected error for amount %f, got success", tc.amount)
			}
		})
	}
}

func TestProcessPayment_DifferentCurrencies(t *testing.T) {
	mastercardProvider := mastercard.GetNewMasterCardPaymentProvider()
	processor := NewPaymentProcessor([]providers.Provider{mastercardProvider})

	testCases := []struct {
		currency string
		amount   float64
		valid    bool
	}{
		{"USD", 100.00, true},
		{"EUR", 85.50, true},
		{"GBP", 75.25, true},
	}

	for _, tc := range testCases {
		t.Run(tc.currency, func(t *testing.T) {
			request := providers.PaymentRequest{
				Mode:        "mastercard",
				Amount:      tc.amount,
				Currency:    tc.currency,
				CardNumber:  "5555555555554444",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "123",
			}

			_, err := processor.ProcessPayment(request)

			if tc.valid && err != nil {
				t.Errorf("Expected success for currency %s, got error: %v", tc.currency, err)
			}

			if !tc.valid && err == nil {
				t.Errorf("Expected error for currency %s, got success", tc.currency)
			}
		})
	}
}
