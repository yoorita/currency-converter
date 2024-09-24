package validations

import (
	"context"

	converter "github.com/yoorita/currency-converter/api"
	"github.com/go-masonry/mortar/interfaces/log"
	"go.uber.org/fx"
)

type (
	CurrencyConverterValidations interface {
		ValidateCurruncyConvertRequest(context.Context, *converter.ConvertRequest) error
	}

	currencyConverterValidationsImplDeps struct {
		fx.In
		Logger log.Logger
	}

	currencyConverterValidationsImpl struct {
		deps currencyConverterValidationsImplDeps
	}
)

func CreateCurrencyConverterValidations(deps currencyConverterValidationsImplDeps) CurrencyConverterValidations {
	return &currencyConverterValidationsImpl{
		deps: deps,
	}
}

func (impl *currencyConverterValidationsImpl) ValidateCurruncyConvertRequest(ctx context.Context, req *converter.ConvertRequest) error {
	return nil
}
