package services

import (
	"context"

	converter "github.com/yoorita/currency-converter/api"
	"github.com/yoorita/currency-converter/app/controllers"
	"github.com/yoorita/currency-converter/app/validations"

	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/interfaces/monitor"
	"go.uber.org/fx"
)

type currencyConverterServiceImplDeps struct {
	fx.In
	Logger log.Logger
	Validations validations.CurrencyConverterValidations
	Controller  controllers.CurrencyConverterController
	Metrics     monitor.Metrics `optional:"true"`
}

type currencyConverterServiceImpl struct {
	converter.UnimplementedCurrencyConverterServer
	deps currencyConverterServiceImplDeps
}

func CreateCurrencyConverterdService(deps currencyConverterServiceImplDeps) converter.CurrencyConverterServer {
	return &currencyConverterServiceImpl{
		deps: deps,
	}
}

func (impl *currencyConverterServiceImpl) Convert(ctx context.Context, req *converter.ConvertRequest) (res *converter.ConvertResponse, err error) {
	err = impl.deps.Validations.ValidateCurruncyConvertRequest(ctx, req)
	if err != nil {
		impl.deps.Logger.WithError(err).WithField("request", req).Error(ctx, "validation failed")
		return nil, err
	}
	return impl.deps.Controller.Convert(ctx, req)
}
