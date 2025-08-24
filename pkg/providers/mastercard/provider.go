package mastercard

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"pgas/pkg/providers"
	"strconv"
	"time"
)

type MasterCardPaymentProvider struct {
	Name string
}

func GetNewMasterCardPaymentProvider() *MasterCardPaymentProvider {
	return &MasterCardPaymentProvider{Name: "mastercard"}
}

func (p *MasterCardPaymentProvider) GetName() string {
	return p.Name
}

func (p *MasterCardPaymentProvider) ValidateRequest(request providers.PaymentRequest) error {

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

func (p *MasterCardPaymentProvider) ProcessPayment(ctx context.Context, request providers.PaymentRequest) (interface{}, interface{}) {
	// Simulate a dummy error response sometimes
	if rand.Float64() < 0.1 {
		errorResponse := map[string]interface{}{
			"error_code": "MC0001",
			"message":    "Insufficient funds",
		}
		return nil, errorResponse
	}

	// Simulate a dummy successful payment response
	successResponse := map[string]interface{}{
		"transaction_id": "TX1234567890",
		"status":         "APPROVED",
		"amount":         strconv.FormatFloat(request.Amount, 'f', -1, 64),
		"currency":       request.Currency,
		"timestamp":      time.Now(),
	}

	return successResponse, nil
}

func (p *MasterCardPaymentProvider) ParseSuccessResponse(response interface{}) (*providers.PaymentResponse, error) {
	data, ok := response.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected map[string]interface{}, got %T", response)
	}

	amountStr, ok := data["amount"].(string)
	if !ok {
		return nil, fmt.Errorf("expected 'amount' field to be a string, got %T", data["amount"])
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to convert 'amount' to float64: %v", err)
	}

	dt, _ := data["timestamp"].(time.Time)

	responseObj := &providers.PaymentResponse{
		Success:       true,
		TransactionID: data["transaction_id"].(string),
		Status:        data["status"].(string),
		Amount:        amount,
		Currency:      data["currency"].(string),
		Date:          &dt,
	}

	return responseObj, nil
}

func (p *MasterCardPaymentProvider) ParseErrorResponse(response interface{}) (*providers.PaymentError, error) {
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
		ErrorCode:    providerError.ErrorCode,
		ErrorMessage: providerError.Message,
	}, nil
}
