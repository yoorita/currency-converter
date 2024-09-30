package data

import (
	"context"
	"errors"
	"strconv"

	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/vtopc/go-monobank"
	"github.com/yoorita/currency-converter/app/clients"
	"go.uber.org/fx"
)

type(
	CurrencyConverterDao interface {
		GetRates(context.Context, string, string) (monobank.Currency, error)
	}

	currencyConverterDaoImplDeps struct {
		fx.In

		Logger    log.Logger
		Monobank clients.MonobankClient
	}

	currencyConverterDaoImpl struct {
		deps currencyConverterDaoImplDeps
	}
)

func CreateCurrencyConverterDao(deps currencyConverterDaoImplDeps) CurrencyConverterDao {
	return &currencyConverterDaoImpl{
		deps: deps,
	}
}

func (impl *currencyConverterDaoImpl) GetRates(ctx context.Context, from, to string) (rates monobank.Currency, err error) {
	monobankRates, err := impl.deps.Monobank.GetCurrencyRates(ctx)
	if err != nil {
		impl.deps.Logger.WithError(err).Error(ctx, "failed to fetch monobank currency rates")
		return
	}
	for _, rate := range monobankRates {
		strA := strconv.Itoa(rate.CurrencyCodeA)
		strB := strconv.Itoa(rate.CurrencyCodeB)

		if (strA == from && strB == to) || (strB == from && strA == to) {
			rates = rate
			return
		}
	}
	err = errors.New("couldn't find match rate for currencies")
	return
}