package controllers

import (
	"context"

	converter "github.com/yoorita/currency-converter/api"

	"github.com/go-masonry/mortar/interfaces/log"
	"go.uber.org/fx"
)

type (
	CurrencyConverterController interface {
		converter.CurrencyConverterServer
	}

	currencyConverterControllerImplDeps struct {
		fx.In
		Logger log.Logger
		Lifecycle fx.Lifecycle
	}

	currencyConverterControllerImpl struct {
		*converter.UnimplementedCurrencyConverterServer
		deps currencyConverterControllerImplDeps
	}
)

func CreateCurrencyConverterController(deps currencyConverterControllerImplDeps) CurrencyConverterController {
	return &currencyConverterControllerImpl{
		deps: deps,
	}
}

func (impl *currencyConverterControllerImpl) Convert(ctx context.Context, req *converter.ConvertRequest) (res *converter.ConvertResponse, err error) {
	impl.deps.Logger.WithError(err).WithField("request", req).WithField("result", res).Info(ctx, "finished conversion")
	return
}
