package clients

import (
	"context"

	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/vtopc/go-monobank"
	"go.uber.org/fx"
)

type (
	MonobankClient interface {
		GetCurrencyRates(context.Context) (monobank.Currencies, error)
	}

	monobankClientImplDeps struct {
		fx.In

		Logger log.Logger
	}

	monobankClientImpl struct {
		deps monobankClientImplDeps
		client monobank.Client
	}
)

func CreateMonobankClient(deps monobankClientImplDeps) MonobankClient {
	return &monobankClientImpl{
		deps: deps,
		client: monobank.NewClient(nil),
	}
}

func (impl *monobankClientImpl) GetCurrencyRates(ctx context.Context) (curr monobank.Currencies, err error) {
	curr, err = impl.client.Currency(ctx)
	impl.deps.Logger.WithError(err).Info(ctx, "got currencies")
	return
}

