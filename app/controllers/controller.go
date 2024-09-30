package controllers

import (
	"context"
	"strconv"

	converter "github.com/yoorita/currency-converter/api"
	"github.com/yoorita/currency-converter/app/clients"
	"github.com/yoorita/currency-converter/app/data"

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
		Sqlite *clients.LazySQLCLient
		Currencies data.CurrencyConverterDao
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
	fromCode, err := impl.deps.Sqlite.Client.GetCurrencyCode(ctx, req.GetCurrencyFrom())
	res = new(converter.ConvertResponse)

	if err != nil {
		impl.deps.Logger.WithError(err).Error(ctx, "failed fetching code for FROM currency")
		return
	}

	toCode, err := impl.deps.Sqlite.Client.GetCurrencyCode(ctx, req.GetCurrencyTo())

	if err != nil {
		impl.deps.Logger.WithError(err).Error(ctx, "failed fetching code for TO currency")
		return
	}

	rate, err := impl.deps.Currencies.GetRates(ctx, fromCode, toCode)

	if err != nil {
		impl.deps.Logger.WithError(err).Error(ctx, "failed to get currancy rates")
		return
	}

	res.Currency = req.GetCurrencyTo()
	
	if strconv.Itoa(rate.CurrencyCodeA) == fromCode {
		res.Amount = req.GetAmountFrom() * float32(rate.RateBuy)
	} else {
		res.Amount = req.GetAmountFrom() / float32(rate.RateSell)
	}

	impl.deps.Logger.WithError(err).WithField("res", res).Info(ctx, "finished conversion")
	return
}
