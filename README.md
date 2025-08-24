# Payment Gateway Adapter System (PGAS)

A robust payment processing system that can work with multiple payment providers. Each provider has different response formats, error structures, and operational characteristics. The system is designed to be easily extensible to support new providers while maintaining high performance and reliability.

## Architecture Overview

The system follows a **modular architecture** with clear separation of concerns:

```
Request → PaymentProcessor → Provider → Response Normalization
```

### Core Components

1. **PaymentRequest**: Contains payment data including provider name
2. **PaymentProcessor**: Main orchestrator that routes requests to appropriate providers
3. **Provider Interface**: Contract that all payment providers must implement
4. **PaymentResponse/PaymentError**: Normalized response formats

### Provider Responsibilities

Each provider must implement:

- **Request Validator**: Provider-specific validation logic
- **Request Processor**: API calls and business logic
- **Success Response Parser**: Convert provider success responses to normalized format
- **Error Response Parser**: Convert provider error responses to normalized format
- **Provider Name**: Unique identifier for the provider

## Features

- **Modular Design**: Clear separation of concerns with provider-specific implementations
- **Provider-Specific Validation**: Each provider can implement its own validation rules
- **Response Normalization**: Convert different provider responses to common internal format
- **Error Handling**: Comprehensive error handling with provider-specific error types
- **Extensible Architecture**: Easy to add new payment providers
- **Thread-Safe**: Supports concurrent payment processing
- **Context Support**: Full context and timeout support for payment processing

## Supported Providers

### Provider A
- **Response Format**: Uses integer amounts (cents)
- **Status Values**: "APPROVED"
- **Error Format**: `{"error_code": "...", "message": "..."}`
- **Validation**: Amount limits, card number validation, expiry validation
- **Special Features**: Simulates 10% random failure rate

### Provider B
- **Response Format**: Uses string amounts with decimal places
- **Status Values**: "SUCCESS"
- **Error Format**: `{"errorType": "...", "reason": "...", "details": {"code": "..."}}`
- **Validation**: Currency restrictions, expiry date validation, amount limits
- **Special Features**: Simulates 15% random failure rate

## Setup and Installation

### Prerequisites
- Go 1.24.3 or higher

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd pgas
```

2. Run the main application:
```bash
go run main.go
```

3. Run tests:
```bash
go test ./pkg/processor/...
```

## Usage Examples

### Basic Payment Processing

```go
package main

import (
    "context"
    "time"
    "pgas/pkg/processor"
)

func main() {
    // Initialize the payment processor
    paymentProcessor := processor.NewPaymentProcessor()
    
    // Create a payment request with provider name included
    request := processor.PaymentRequest{
        ProviderName: "provider_a",  // Specify which provider to use
        Amount:       100.00,
        Currency:     "USD",
        CardNumber:   "4111111111111111",
        ExpiryMonth:  "12",
        ExpiryYear:   "2025",
        CVV:          "123",
    }
    
    // Process payment with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    response, err := paymentProcessor.ProcessPayment(ctx, request)
    if err != nil {
        // Handle payment-specific errors
        if paymentError, ok := err.(*processor.PaymentError); ok {
            fmt.Printf("Payment failed: %s - %s\n", paymentError.Code, paymentError.Message)
        }
        return
    }
    
    fmt.Printf("Payment successful! Transaction ID: %s\n", response.TransactionID)
}
```



```

# Adding New Payment Providers

This guide explains how to add new payment providers to the Payment Gateway Adapter System (PGAS).

## Overview

The system is designed to be easily extensible. Adding a new provider involves implementing the `Provider` interface and following a consistent pattern.

## Provider Interface

All payment providers must implement the `Provider` interface defined in `pkg/providers/types.go`:

```go
type Provider interface {
    GetName() string
    ValidateRequest(request PaymentRequest) error
    ProcessPayment(ctx context.Context, request PaymentRequest) (interface{}, interface{})
    ParseSuccessResponse(response interface{}) (*PaymentResponse, error)
    ParseErrorResponse(response interface{}) (*PaymentError, error)
}
```

## Step-by-Step Guide

### Step 1: Create Provider Directory

Create a new directory for your provider under `pkg/providers/`:

```bash
mkdir pkg/providers/your_provider_name
cd pkg/providers/your_provider_name
```

### Step 2: Create Provider Implementation

Create a file named `provider.go` in your provider directory:

```go
package your_provider_name

import (
    "context"
    "errors"
    "fmt"
    "strconv"
    "time"
    
    "pgas/pkg/providers"
)

// Provider-specific response structures
type YourProviderResponse struct {
    TransactionID string  `json:"transaction_id"`
    Status        string  `json:"status"`
    Amount        float64 `json:"amount"`
    Currency      string  `json:"currency"`
    Timestamp     string  `json:"timestamp"`
}

type YourProviderError struct {
    ErrorCode    string `json:"error_code"`
    ErrorMessage string `json:"error_message"`
    Details      string `json:"details,omitempty"`
}

// NewYourProvider creates a new instance of YourProvider
func NewYourProvider(apiKey, baseURL string) *YourProvider {
    return &YourProvider{
        apiKey:  apiKey,
        baseURL: baseURL,
        timeout: 30 * time.Second,
    }
}

// GetName returns the provider name
func (p *YourProvider) GetName() string {
    return "your_provider_name"
}

// ValidateRequest validates the payment request
func (p *YourProvider) ValidateRequest(request providers.PaymentRequest) error {
    return nil
}

// ProcessPayment processes the payment request
func (p *YourProvider) ProcessPayment(ctx context.Context, request providers.PaymentRequest) (interface{}, interface{}) {
   
    return response, nil
}

// ParseSuccessResponse converts provider response to normalized format
func (p *YourProvider) ParseSuccessResponse(response interface{}) (*providers.PaymentResponse, error) {
    return normalizedResponse, nil
}

// ParseErrorResponse converts provider error to normalized format
func (p *YourProvider) ParseErrorResponse(response interface{}) (*providers.PaymentError, error) {
    return normalizedError, nil
}
```

### Step 4: Create Tests

Create a test file `provider_test.go` in your provider directory and write all tests in it

### Step 5: Register Your Provider

Update the main application to include your new provider:

```go
// In main.go
import (
    "pgas/pkg/providers/your_provider_name"
)

func main() {
    // Initialize payment providers
    yourProvider := your_provider_name.GetNewYourProviderPaymentProvider()
    
    // Add to processor
    paymentProcessor := processor.NewPaymentProcessor([]providers.Provider{
        mastercardProvider, 
        visaProvider, 
        yourProvider, // Add your provider here
    })
    
}
```

### Step 6: Update Integration Tests

Add your provider to the integration tests in `pkg/processor/integration_test.go`:
## Best Practices

### 1. Validation
- Implement comprehensive validation for all input fields
- Check for supported currencies, valid card numbers, expiry dates, etc.
- Return descriptive error messages

### 2. Error Handling
- Handle context cancellation
- Implement proper error types and codes
- Provide meaningful error messages

### 3. Response Parsing
- Ensure consistent response format
- Handle different response structures
- Validate response data before parsing

### 4. Testing
- Write comprehensive unit tests
- Test both success and failure scenarios
- Include edge cases and boundary conditions
- Test concurrent access if applicable

### 5. Configuration
- Use environment variables for sensitive data (API keys, URLs)
- Implement proper timeout handling
- Add retry logic for transient failures
