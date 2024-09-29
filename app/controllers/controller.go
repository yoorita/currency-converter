package controllers

import (
	"context"

	converter "github.com/yoorita/currency-converter/api"
	"github.com/yoorita/currency-converter/app/clients"

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
		Monobank clients.MonobankClient
		Sqlite *clients.LazySQLCLient
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
	// impl.deps.Monobank.GetCurrencyRates(ctx)
	fromCode, err := impl.deps.Sqlite.Client.GetCurrencyCode(ctx, req.GetCurrencyFrom())
	impl.deps.Logger.WithError(err).WithField("from", fromCode).Info(ctx, "finished conversion")
	return
}
