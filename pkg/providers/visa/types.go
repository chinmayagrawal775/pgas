package visa

type PaymentRequest struct {
	// request format for visa
}

// success response format for visa
type PaymentResponse struct {
	PaymentID string `json:"payment_id"`
	State     string `json:"state"`
	Value     struct {
		Amount       string `json:"amount"`
		CurrencyCode string `json:"currency_code"`
	} `json:"value"`
	ProcessedAt int64 `json:"processed_at"`
}

// error response format for visa
type PaymentError struct {
	ErrorType string `json:"error_type"`
	Reason    string `json:"reason"`
	Details   struct {
		Code string `json:"code"`
	} `json:"details"`
}
