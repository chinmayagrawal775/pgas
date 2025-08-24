package mastercard

import "time"

// request format for mastercard
type PaymentRequest struct {
}

// success response format for mastercard
type PaymentResponse struct {
	TransactionID string    `json:"transaction_id"`
	Status        string    `json:"status"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Timestamp     time.Time `json:"timestamp"` // eg: "2024-01-15T10:30:00Z"
}

// error response format for mastercard
type PaymentError struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}
