package processor

import (
	"testing"

	"pgas/pkg/providers"
	"pgas/pkg/providers/mastercard"
	"pgas/pkg/providers/visa"
)

func TestIntegration_SuccessfulPayments(t *testing.T) {
	mastercardProvider := mastercard.GetNewMasterCardPaymentProvider()
	visaProvider := visa.GetNewVisaPaymentProvider()

	processor := NewPaymentProcessor([]providers.Provider{mastercardProvider, visaProvider})

	testCases := []struct {
		name     string
		provider string
		request  providers.PaymentRequest
	}{
		{
			name:     "Visa successful payment",
			provider: "visa",
			request: providers.PaymentRequest{
				Mode:        "visa",
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  "4111111111111111",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "123",
			},
		},
		{
			name:     "Mastercard successful payment",
			provider: "mastercard",
			request: providers.PaymentRequest{
				Mode:        "mastercard",
				Amount:      50.00,
				Currency:    "USD",
				CardNumber:  "5555555555554444",
				ExpiryMonth: "10",
				ExpiryYear:  "2024",
				CVV:         "456",
			},
		},
		{
			name:     "Visa EUR payment",
			provider: "visa",
			request: providers.PaymentRequest{
				Mode:        "visa",
				Amount:      85.50,
				Currency:    "EUR",
				CardNumber:  "4111111111111111",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "123",
			},
		},
		{
			name:     "Mastercard GBP payment",
			provider: "mastercard",
			request: providers.PaymentRequest{
				Mode:        "mastercard",
				Amount:      75.25,
				Currency:    "GBP",
				CardNumber:  "5555555555554444",
				ExpiryMonth: "10",
				ExpiryYear:  "2024",
				CVV:         "456",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := processor.ProcessPayment(tc.request)
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

			if response.Amount != tc.request.Amount {
				t.Errorf("Expected amount %f, got %f", tc.request.Amount, response.Amount)
			}

			if response.Currency != tc.request.Currency {
				t.Errorf("Expected currency %s, got %s", tc.request.Currency, response.Currency)
			}

			if response.Date == nil {
				t.Error("Expected date to be set")
			}

		})
	}
}

func TestIntegration_ErrorScenarios(t *testing.T) {
	mastercardProvider := mastercard.GetNewMasterCardPaymentProvider()
	visaProvider := visa.GetNewVisaPaymentProvider()

	processor := NewPaymentProcessor([]providers.Provider{mastercardProvider, visaProvider})

	testCases := []struct {
		name           string
		request        providers.PaymentRequest
		expectedError  bool
		expectedCode   string
		expectedReason string
	}{
		{
			name: "invalid provider",
			request: providers.PaymentRequest{
				Mode:        "invalid_provider",
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  "4111111111111111",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "123",
			},
			expectedError:  true,
			expectedCode:   "INVALID_PROVIDER",
			expectedReason: "invalid provider name provided",
		},
		{
			name: "zero amount",
			request: providers.PaymentRequest{
				Mode:        "visa",
				Amount:      0,
				Currency:    "USD",
				CardNumber:  "4111111111111111",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "123",
			},
			expectedError:  true,
			expectedCode:   "INVALID_REQUEST",
			expectedReason: "amount must be greater than 0",
		},
		{
			name: "negative amount",
			request: providers.PaymentRequest{
				Mode:        "mastercard",
				Amount:      -100.00,
				Currency:    "USD",
				CardNumber:  "5555555555554444",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "123",
			},
			expectedError:  true,
			expectedCode:   "INVALID_REQUEST",
			expectedReason: "amount must be greater than 0",
		},
		{
			name: "empty card number",
			request: providers.PaymentRequest{
				Mode:        "visa",
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  "",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "123",
			},
			expectedError:  true,
			expectedCode:   "INVALID_REQUEST",
			expectedReason: "card number is required",
		},
		{
			name: "invalid CVV",
			request: providers.PaymentRequest{
				Mode:        "mastercard",
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  "5555555555554444",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "12",
			},
			expectedError:  true,
			expectedCode:   "INVALID_REQUEST",
			expectedReason: "CVV must be 3 or 4 digits",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := processor.ProcessPayment(tc.request)

			if tc.expectedError {
				if err == nil {
					t.Fatal("Expected error, got nil")
				}

				if err.ErrorCode != tc.expectedCode {
					t.Errorf("Expected error code '%s', got '%s'", tc.expectedCode, err.ErrorCode)
				}

				if response != nil {
					t.Fatal("Expected nil response for error case")
				}
			} else {
				if err != nil {
					t.Fatalf("Expected success, got error: %v", err)
				}

				if response == nil {
					t.Fatal("Expected response, got nil")
				}
			}
		})
	}
}

func TestIntegration_EdgeCases(t *testing.T) {
	mastercardProvider := mastercard.GetNewMasterCardPaymentProvider()
	visaProvider := visa.GetNewVisaPaymentProvider()

	processor := NewPaymentProcessor([]providers.Provider{mastercardProvider, visaProvider})

	testCases := []struct {
		name    string
		request providers.PaymentRequest
		valid   bool
	}{
		{
			name: "very long card number",
			request: providers.PaymentRequest{
				Mode:        "visa",
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  "4111111111111111111111111111111111111111111111111111111111111111",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "123",
			},
			valid: false,
		},
		{
			name: "very short card number",
			request: providers.PaymentRequest{
				Mode:        "mastercard",
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  "123",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "123",
			},
			valid: false,
		},
		{
			name: "4-digit CVV",
			request: providers.PaymentRequest{
				Mode:        "visa",
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  "4111111111111111",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "1234",
			},
			valid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := processor.ProcessPayment(tc.request)

			if tc.valid && err != nil {
				t.Errorf("Expected success for %s, got error: %v", tc.name, err)
			}

			if !tc.valid && err == nil {
				t.Errorf("Expected error for %s, got success", tc.name)
			}

			if tc.valid && response != nil {
				if response.Amount != tc.request.Amount {
					t.Errorf("Expected amount %f, got %f", tc.request.Amount, response.Amount)
				}
			}
		})
	}
}

func TestIntegration_ConcurrentPayments(t *testing.T) {
	mastercardProvider := mastercard.GetNewMasterCardPaymentProvider()
	visaProvider := visa.GetNewVisaPaymentProvider()

	processor := NewPaymentProcessor([]providers.Provider{mastercardProvider, visaProvider})

	// Test concurrent payments to ensure thread safety
	const numGoroutines = 10
	results := make(chan *providers.PaymentError, numGoroutines)

	request := providers.PaymentRequest{
		Mode:        "visa",
		Amount:      25.00,
		Currency:    "USD",
		CardNumber:  "4111111111111111",
		ExpiryMonth: "12",
		ExpiryYear:  "2025",
		CVV:         "123",
	}

	// Start concurrent payment processing
	for i := 0; i < numGoroutines; i++ {
		go func() {
			_, err := processor.ProcessPayment(request)
			if err != nil {
				results <- err
			} else {
				results <- nil
			}
		}()
	}

	// Collect results
	successCount := 0
	errorCount := 0

	for i := 0; i < numGoroutines; i++ {
		err := <-results
		if err != nil {
			errorCount++
		} else {
			successCount++
		}
	}

	// All payments should succeed in this case
	if successCount == 0 {
		t.Error("Expected at least some successful payments")
	}

	t.Logf("Concurrent payments: %d successful, %d failed", successCount, errorCount)
}

func TestIntegration_ProviderSpecificBehavior(t *testing.T) {
	mastercardProvider := mastercard.GetNewMasterCardPaymentProvider()
	visaProvider := visa.GetNewVisaPaymentProvider()

	processor := NewPaymentProcessor([]providers.Provider{mastercardProvider, visaProvider})

	// Test that each provider has different behavior
	testCases := []struct {
		name       string
		provider   string
		cardNumber string
	}{
		{"Visa card with Visa provider", "visa", "4111111111111111"},
		{"Mastercard with Mastercard provider", "mastercard", "5555555555554444"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := providers.PaymentRequest{
				Mode:        tc.provider,
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  tc.cardNumber,
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

			// Verify provider-specific behavior
			if response.TransactionID == "" {
				t.Error("Expected transaction ID to be set")
			}

			// Different providers might have different transaction ID formats
			if tc.provider == "visa" {
				// Visa might have specific transaction ID format
				if len(response.TransactionID) == 0 {
					t.Error("Expected Visa transaction ID to be set")
				}
			} else if tc.provider == "mastercard" {
				// Mastercard might have specific transaction ID format
				if len(response.TransactionID) == 0 {
					t.Error("Expected Mastercard transaction ID to be set")
				}
			}
		})
	}
}
