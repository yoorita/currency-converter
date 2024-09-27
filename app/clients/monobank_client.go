package clients

import (
	"context"

	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/vtopc/go-monobank"
	"go.uber.org/fx"
)

type (
	MonobankClient interface {
		GetCurrencyRates(context.Context) error
	}

	monobankClientImplDeps struct {
		fx.In

		Logger log.Logger
		Config cfg.Config
		Lifecycle fx.Lifecycle
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

func (impl *monobankClientImpl) GetCurrencyRates(ctx context.Context) (err error) {
	currency, err := impl.client.Currency(ctx)
	impl.deps.Logger.WithError(err).WithField("result", currency).Info(ctx, "finished")
	return
}

