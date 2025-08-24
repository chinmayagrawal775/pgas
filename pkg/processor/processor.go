package processor

import (
	"context"
	"errors"
	"pgas/pkg/providers"
)

type PaymentProcessor struct {
	providers map[string]providers.Provider
}

func NewPaymentProcessor(paymentProviders []providers.Provider) *PaymentProcessor {
	newProvider := &PaymentProcessor{
		providers: make(map[string]providers.Provider),
	}

	newProvider.registerProviders(paymentProviders)

	return newProvider
}

func (p *PaymentProcessor) registerProviders(providers []providers.Provider) {
	for _, provider := range providers {
		p.providers[provider.GetName()] = provider
	}
}

func (p *PaymentProcessor) getProvider(requiredProvider string) (providers.Provider, error) {
	pr := p.providers[requiredProvider]
	if pr == nil {
		return nil, errors.New("invalid provider name provided: '" + requiredProvider + "'")
	}

	return pr, nil
}

func (p *PaymentProcessor) ProcessPayment(paymentReqest providers.PaymentRequest) (*providers.PaymentResponse, *providers.PaymentError) {

	paymentProvider, err := p.getProvider(paymentReqest.Mode)
	if err != nil {
		return nil, &providers.PaymentError{
			Success:      false,
			ErrorCode:    "INVALID_PROVIDER",
			ErrorMessage: err.Error(),
		}
	}

	validationError := paymentProvider.ValidateRequest(paymentReqest)
	if validationError != nil {
		return nil, &providers.PaymentError{
			Success:      false,
			ErrorCode:    "INVALID_REQUEST",
			ErrorMessage: validationError.Error(),
		}
	}

	ctx := context.Background()

	processResponse, processError := paymentProvider.ProcessPayment(ctx, paymentReqest)

	if processError != nil {

		parseErrorRes, parseErroErr := paymentProvider.ParseErrorResponse(processError)
		if parseErroErr != nil {
			return nil, &providers.PaymentError{
				Success:      false,
				ErrorCode:    "PROCESSING_ERROR",
				ErrorMessage: parseErroErr.Error(),
			}
		}

		return nil, parseErrorRes

	}

	successResponse, successParseError := paymentProvider.ParseSuccessResponse(processResponse)
	if successParseError != nil {
		return nil, &providers.PaymentError{
			Success:      false,
			ErrorCode:    "PARSING_ERROR",
			ErrorMessage: successParseError.Error(),
		}
	}

	return successResponse, nil
}
