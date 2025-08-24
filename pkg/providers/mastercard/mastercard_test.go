package mastercard

import (
	"testing"
	"time"

	"pgas/pkg/providers"
)

func TestGetNewMasterCardPaymentProvider(t *testing.T) {
	provider := GetNewMasterCardPaymentProvider()
	if provider == nil {
		t.Fatal("Expected provider to be created")
	}

	if provider.GetName() != "mastercard" {
		t.Errorf("Expected provider name 'mastercard', got: %s", provider.GetName())
	}
}

func TestMastercardProvider_ValidateRequest(t *testing.T) {
	provider := GetNewMasterCardPaymentProvider()

	testCases := []struct {
		name    string
		request providers.PaymentRequest
		valid   bool
	}{
		{
			name: "valid request",
			request: providers.PaymentRequest{
				Mode:        "mastercard",
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  "5555555555554444",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "123",
			},
			valid: true,
		},
		{
			name: "zero amount",
			request: providers.PaymentRequest{
				Mode:        "mastercard",
				Amount:      0,
				Currency:    "USD",
				CardNumber:  "5555555555554444",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "123",
			},
			valid: false,
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
			valid: false,
		},
		{
			name: "empty currency",
			request: providers.PaymentRequest{
				Mode:        "mastercard",
				Amount:      100.00,
				Currency:    "",
				CardNumber:  "5555555555554444",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "123",
			},
			valid: false,
		},
		{
			name: "empty card number",
			request: providers.PaymentRequest{
				Mode:        "mastercard",
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  "",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "123",
			},
			valid: false,
		},
		{
			name: "invalid card number length",
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
			name: "empty expiry month",
			request: providers.PaymentRequest{
				Mode:        "mastercard",
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  "5555555555554444",
				ExpiryMonth: "",
				ExpiryYear:  "2025",
				CVV:         "123",
			},
			valid: false,
		},
		{
			name: "empty expiry year",
			request: providers.PaymentRequest{
				Mode:        "mastercard",
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  "5555555555554444",
				ExpiryMonth: "12",
				ExpiryYear:  "",
				CVV:         "123",
			},
			valid: false,
		},
		{
			name: "empty CVV",
			request: providers.PaymentRequest{
				Mode:        "mastercard",
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  "5555555555554444",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "",
			},
			valid: false,
		},
		{
			name: "invalid CVV length",
			request: providers.PaymentRequest{
				Mode:        "mastercard",
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  "5555555555554444",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "12",
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := provider.ValidateRequest(tc.request)
			if tc.valid && err != nil {
				t.Errorf("Expected valid request, got error: %v", err)
			}
			if !tc.valid && err == nil {
				t.Errorf("Expected invalid request, got no error")
			}
		})
	}
}

// func TestMastercardProvider_ProcessPayment(t *testing.T) {
// 	provider := GetNewMasterCardPaymentProvider()

// 	request := providers.PaymentRequest{
// 		Mode:        "mastercard",
// 		Amount:      100.00,
// 		Currency:    "USD",
// 		CardNumber:  "5555555555554444",
// 		ExpiryMonth: "12",
// 		ExpiryYear:  "2025",
// 		CVV:         "123",
// 	}

// 	ctx := context.Background()
// 	response, err := provider.ProcessPayment(ctx, request)

// 	if err != nil {
// 		t.Fatalf("Expected successful processing, got error: %v", err)
// 	}

// 	if response == nil {
// 		t.Fatal("Expected response, got nil")
// 	}

// 	// Verify response structure
// 	mastercardResponse, ok := response.(PaymentResponse)
// 	if !ok {
// 		t.Fatal("Expected MastercardResponse type")
// 	}

// 	if mastercardResponse.TransactionID == "" {
// 		t.Error("Expected transaction ID to be set")
// 	}

// 	if mastercardResponse.Status == "" {
// 		t.Error("Expected status to be set")
// 	}

// 	if mastercardResponse.Amount != request.Amount {
// 		t.Errorf("Expected amount %f, got %f", request.Amount, mastercardResponse.Amount)
// 	}

// 	if mastercardResponse.Currency != request.Currency {
// 		t.Errorf("Expected currency %s, got %s", request.Currency, mastercardResponse.Currency)
// 	}
// }

func TestMastercardProvider_ParseSuccessResponse(t *testing.T) {
	provider := GetNewMasterCardPaymentProvider()

	mastercardResponse := map[string]interface{}{
		"transaction_id": "TX1234567890",
		"status":         "APPROVED",
		"amount":         "24.44",
		"currency":       "USD",
		"timestamp":      time.Now(),
	}

	response, err := provider.ParseSuccessResponse(mastercardResponse)
	if err != nil {
		t.Fatalf("Expected successful parsing, got error: %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}

	if response.TransactionID != "TX1234567890" {
		t.Errorf("Expected transaction ID %s, got %s", "TX1234567890", response.TransactionID)
	}

	if response.Status != "APPROVED" {
		t.Errorf("Expected status %s, got %s", "APPROVED", response.Status)
	}

	if response.Amount != 24.44 {
		t.Errorf("Expected amount %f, got %f", 24.44, response.Amount)
	}

	if response.Currency != "USD" {
		t.Errorf("Expected currency %s, got %s", "USD", response.Currency)
	}

	if response.Date == nil {
		t.Error("Expected date to be set")
	}
}

func TestMastercardProvider_ParseErrorResponse(t *testing.T) {
	provider := GetNewMasterCardPaymentProvider()

	mastercardError := map[string]interface{}{
		"error_code": "MC0001",
		"message":    "Insufficient funds",
	}

	errorResponse, err := provider.ParseErrorResponse(mastercardError)
	if err != nil {
		t.Fatalf("Expected successful error parsing, got error: %v", err)
	}

	if errorResponse == nil {
		t.Fatal("Expected error response, got nil")
	}

	if errorResponse.Success {
		t.Error("Expected success to be false")
	}

	if errorResponse.ErrorCode != "MC0001" {
		t.Errorf("Expected error code %s, got %s", "MC0001", errorResponse.ErrorCode)
	}

	if errorResponse.ErrorMessage != "Insufficient funds" {
		t.Errorf("Expected error message %s, got %s", "Insufficient funds", errorResponse.ErrorMessage)
	}
}

func TestMastercardProvider_EdgeCases(t *testing.T) {
	provider := GetNewMasterCardPaymentProvider()

	testCases := []struct {
		name    string
		request providers.PaymentRequest
		valid   bool
	}{
		{
			name: "minimum amount",
			request: providers.PaymentRequest{
				Mode:        "mastercard",
				Amount:      0.01,
				Currency:    "USD",
				CardNumber:  "5555555555554444",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "123",
			},
			valid: true,
		},
		{
			name: "large amount",
			request: providers.PaymentRequest{
				Mode:        "mastercard",
				Amount:      999999.99,
				Currency:    "USD",
				CardNumber:  "5555555555554444",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "123",
			},
			valid: true,
		},
		{
			name: "4-digit CVV",
			request: providers.PaymentRequest{
				Mode:        "mastercard",
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  "5555555555554444",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "1234",
			},
			valid: true,
		},
		{
			name: "future expiry date",
			request: providers.PaymentRequest{
				Mode:        "mastercard",
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  "5555555555554444",
				ExpiryMonth: "12",
				ExpiryYear:  "2030",
				CVV:         "123",
			},
			valid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := provider.ValidateRequest(tc.request)
			if tc.valid && err != nil {
				t.Errorf("Expected valid request for %s, got error: %v", tc.name, err)
			}
			if !tc.valid && err == nil {
				t.Errorf("Expected invalid request for %s, got no error", tc.name)
			}
		})
	}
}
