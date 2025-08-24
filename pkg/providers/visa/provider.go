package visa

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand/v2"
	"pgas/pkg/providers"
	"strconv"
	"time"
)

type VisaPaymentProvider struct {
	Name string
}

func GetNewVisaPaymentProvider() *VisaPaymentProvider {
	return &VisaPaymentProvider{Name: "visa"}
}

func (p *VisaPaymentProvider) GetName() string {
	return p.Name
}

func (p *VisaPaymentProvider) ValidateRequest(request providers.PaymentRequest) error {

	if request.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	if request.Amount > 1000000 {
		return errors.New("amount exceeds maximum limit of 1,000,000")
	}

	if request.Currency == "" {
		return errors.New("currency is required")
	}

	if request.CardNumber == "" {
		return errors.New("card number is required")
	}

	if len(request.CardNumber) < 13 || len(request.CardNumber) > 19 {
		return errors.New("card number must be between 13 and 19 digits")
	}

	if request.ExpiryMonth == "" || request.ExpiryYear == "" {
		return errors.New("expiry month and year are required")
	}

	if request.CVV == "" {
		return errors.New("CVV is required")
	}

	if len(request.CVV) < 3 || len(request.CVV) > 4 {
		return errors.New("CVV must be 3 or 4 digits")
	}

	return nil
}

func (p *VisaPaymentProvider) ProcessPayment(ctx context.Context, request providers.PaymentRequest) (interface{}, interface{}) {

	// Simulate a dummy error response sometimes
	if rand.Float64() < 0.1 {
		errorResponse := map[string]interface{}{
			"error_type": "PAYMENT_FAILED",
			"reason":     "Card declined",
			"details": map[string]interface{}{
				"code": "EE000011",
			},
		}
		return nil, errorResponse
	}

	// Simulate a dummy successful payment response
	successResponse := map[string]interface{}{
		"payment_id": "PPAAYY--778899--XXYYZZ",
		"state":      "SUCCESS",
		"value": map[string]interface{}{
			"amount":        strconv.FormatFloat(request.Amount, 'f', -1, 64),
			"currency_code": request.Currency,
		},
		"processed_at": 1677587921,
	}

	return successResponse, nil
}

func (p *VisaPaymentProvider) ParseSuccessResponse(response interface{}) (*providers.PaymentResponse, error) {
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return nil, errors.New("error marshalling response")
	}

	var providerResponse PaymentResponse
	err = json.Unmarshal(responseJSON, &providerResponse)
	if err != nil {
		return nil, errors.New("invalid response type")
	}

	parsedAmount, _ := strconv.ParseFloat(providerResponse.Value.Amount, 64)
	parsedTime := time.Unix(providerResponse.ProcessedAt, 0)

	return &providers.PaymentResponse{
		Success:       true,
		TransactionID: providerResponse.PaymentID,
		Status:        providerResponse.State,
		Amount:        parsedAmount,
		Currency:      providerResponse.Value.CurrencyCode,
		Date:          &parsedTime,
	}, nil
}

func (p *VisaPaymentProvider) ParseErrorResponse(response interface{}) (*providers.PaymentError, error) {
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return nil, errors.New("error marshalling error response")
	}

	var providerError PaymentError
	err = json.Unmarshal(responseJSON, &providerError)
	if err != nil {
		return nil, errors.New("invalid response error type")
	}

	return &providers.PaymentError{
		Success:      false,
		ErrorCode:    providerError.Details.Code,
		ErrorMessage: "ErrorType:" + providerError.ErrorType + " :: ErrorReason: " + providerError.Reason,
	}, nil
}
