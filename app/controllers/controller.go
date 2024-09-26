package controllers

import (
	"context"

	converter "github.com/yoorita/currency-converter/api"

	"github.com/go-masonry/mortar/interfaces/log"
	"go.uber.org/fx"
)

type CurrencyConverterController interface {
	converter.CurrencyConverterServer
}

type currencyConverterControllerImplDeps struct {
	fx.In
	Logger log.Logger
}

type currencyConverterControllerImpl struct {
	*converter.UnimplementedCurrencyConverterServer
	deps currencyConverterControllerImplDeps
}

func CreateCurrencyConverterController(deps currencyConverterControllerImplDeps) CurrencyConverterController {
	return &currencyConverterControllerImpl{
		deps: deps,
	}
}

func (w *currencyConverterControllerImpl) Convert(ctx context.Context, req *converter.ConvertRequest) (res *converter.ConvertResponse, err error) {
	return &converter.ConvertResponse{}, nil
}
