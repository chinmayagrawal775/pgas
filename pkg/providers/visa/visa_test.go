package visa

import (
	"testing"

	"pgas/pkg/providers"
)

func TestGetNewVisaPaymentProvider(t *testing.T) {
	provider := GetNewVisaPaymentProvider()
	if provider == nil {
		t.Fatal("Expected provider to be created")
	}

	if provider.GetName() != "visa" {
		t.Errorf("Expected provider name 'visa', got: %s", provider.GetName())
	}
}

func TestVisaProvider_ValidateRequest(t *testing.T) {
	provider := GetNewVisaPaymentProvider()

	testCases := []struct {
		name    string
		request providers.PaymentRequest
		valid   bool
	}{
		{
			name: "valid request",
			request: providers.PaymentRequest{
				Mode:        "visa",
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  "4111111111111111",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "123",
			},
			valid: true,
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
			valid: false,
		},
		{
			name: "negative amount",
			request: providers.PaymentRequest{
				Mode:        "visa",
				Amount:      -100.00,
				Currency:    "USD",
				CardNumber:  "4111111111111111",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "123",
			},
			valid: false,
		},
		{
			name: "empty currency",
			request: providers.PaymentRequest{
				Mode:        "visa",
				Amount:      100.00,
				Currency:    "",
				CardNumber:  "4111111111111111",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "123",
			},
			valid: false,
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
			valid: false,
		},
		{
			name: "invalid card number length",
			request: providers.PaymentRequest{
				Mode:        "visa",
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
				Mode:        "visa",
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  "4111111111111111",
				ExpiryMonth: "",
				ExpiryYear:  "2025",
				CVV:         "123",
			},
			valid: false,
		},
		{
			name: "empty expiry year",
			request: providers.PaymentRequest{
				Mode:        "visa",
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  "4111111111111111",
				ExpiryMonth: "12",
				ExpiryYear:  "",
				CVV:         "123",
			},
			valid: false,
		},
		{
			name: "empty CVV",
			request: providers.PaymentRequest{
				Mode:        "visa",
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  "4111111111111111",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "",
			},
			valid: false,
		},
		{
			name: "invalid CVV length",
			request: providers.PaymentRequest{
				Mode:        "visa",
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  "4111111111111111",
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

// func TestVisaProvider_ProcessPayment(t *testing.T) {
// 	provider := GetNewVisaPaymentProvider()

// 	request := providers.PaymentRequest{
// 		Mode:        "visa",
// 		Amount:      100.00,
// 		Currency:    "USD",
// 		CardNumber:  "4111111111111111",
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
// 	visaResponse, ok := response.(PaymentResponse)
// 	if !ok {
// 		t.Fatal("Expected VisaResponse type")
// 	}

// 	if visaResponse.TransactionID == "" {
// 		t.Error("Expected transaction ID to be set")
// 	}

// 	if visaResponse.Status == "" {
// 		t.Error("Expected status to be set")
// 	}

// 	if visaResponse.Amount != request.Amount {
// 		t.Errorf("Expected amount %f, got %f", request.Amount, visaResponse.Amount)
// 	}

// 	if visaResponse.Currency != request.Currency {
// 		t.Errorf("Expected currency %s, got %s", request.Currency, visaResponse.Currency)
// 	}
// }

func TestVisaProvider_ParseSuccessResponse(t *testing.T) {
	provider := GetNewVisaPaymentProvider()

	visaResponse := map[string]interface{}{
		"payment_id": "PPAAYY--778899--XXYYZZ",
		"state":      "SUCCESS",
		"value": map[string]interface{}{
			"amount":        "1000.00",
			"currency_code": "USD",
		},
		"processed_at": 1677587921,
	}

	response, err := provider.ParseSuccessResponse(visaResponse)
	if err != nil {
		t.Fatalf("Expected successful parsing, got error: %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}

	if response.TransactionID != "PPAAYY--778899--XXYYZZ" {
		t.Errorf("Expected transaction ID %s, got %s", "PPAAYY--778899--XXYYZZ", response.TransactionID)
	}

	if response.Status != "SUCCESS" {
		t.Errorf("Expected status %s, got %s", "SUCCESS", response.Status)
	}

	if response.Currency != "USD" {
		t.Errorf("Expected currency %s, got %s", "USD", response.Currency)
	}

	if response.Amount != 1000.00 {
		t.Errorf("Expected amount %f, got %f", 1000.00, response.Amount)
	}

	if response.Date == nil {
		t.Error("Expected date to be set")
	}
}

func TestVisaProvider_ParseErrorResponse(t *testing.T) {
	provider := GetNewVisaPaymentProvider()

	visaError := map[string]interface{}{
		"error_type": "PAYMENT_FAILED",
		"reason":     "Card declined",
		"details": map[string]interface{}{
			"code": "EE000011",
		},
	}

	errorResponse, err := provider.ParseErrorResponse(visaError)
	if err != nil {
		t.Fatalf("Expected successful error parsing, got error: %v", err)
	}

	if errorResponse == nil {
		t.Fatal("Expected error response, got nil")
	}

	if errorResponse.Success {
		t.Error("Expected success to be false")
	}

	if errorResponse.ErrorCode != "EE000011" {
		t.Errorf("Expected error code %s, got %s", "EE000011", errorResponse.ErrorCode)
	}

	if errorResponse.ErrorMessage != "ErrorType:PAYMENT_FAILED :: ErrorReason: Card declined" {
		t.Errorf("Expected error message %s, got %s", "ErrorType:PAYMENT_FAILED :: ErrorReason: Card declined", errorResponse.ErrorMessage)
	}
}

func TestVisaProvider_EdgeCases(t *testing.T) {
	provider := GetNewVisaPaymentProvider()

	testCases := []struct {
		name    string
		request providers.PaymentRequest
		valid   bool
	}{
		{
			name: "minimum amount",
			request: providers.PaymentRequest{
				Mode:        "visa",
				Amount:      0.01,
				Currency:    "USD",
				CardNumber:  "4111111111111111",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "123",
			},
			valid: true,
		},
		{
			name: "large amount",
			request: providers.PaymentRequest{
				Mode:        "visa",
				Amount:      999999.99,
				Currency:    "USD",
				CardNumber:  "4111111111111111",
				ExpiryMonth: "12",
				ExpiryYear:  "2025",
				CVV:         "123",
			},
			valid: true,
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
		{
			name: "future expiry date",
			request: providers.PaymentRequest{
				Mode:        "visa",
				Amount:      100.00,
				Currency:    "USD",
				CardNumber:  "4111111111111111",
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
