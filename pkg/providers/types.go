package providers

import (
	"context"
	"time"
)

// normalized request format for internal/user purpose
type PaymentRequest struct {
	Mode        string  `json:"mode"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	CardNumber  string  `json:"card_number"`
	ExpiryMonth string  `json:"expiry_month"`
	ExpiryYear  string  `json:"expiry_year"`
	CVV         string  `json:"cvv"`
}

// normalized success response format for internal/user purpose
type PaymentResponse struct {
	Success       bool       `json:"success"`
	TransactionID string     `json:"transaction_id"`
	Status        string     `json:"status"`
	Amount        float64    `json:"amount,omitempty"`
	Currency      string     `json:"currency,omitempty"`
	Date          *time.Time `json:"date,omitempty"`
}

// normalized error response format for internal/user purpose
type PaymentError struct {
	Success      bool   `json:"success"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

type Provider interface {
	GetName() string
	ValidateRequest(request PaymentRequest) error
	ProcessPayment(ctx context.Context, request PaymentRequest) (interface{}, interface{})
	ParseSuccessResponse(response interface{}) (*PaymentResponse, error)
	ParseErrorResponse(response interface{}) (*PaymentError, error)
}
