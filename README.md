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

### Processing with Multiple Providers

```go
func processWithMultipleProviders(processor *processor.PaymentProcessor, amount float64) {
    providers := []string{"provider_a", "provider_b"}
    
    for _, providerName := range providers {
        re