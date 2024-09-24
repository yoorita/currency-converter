package services

import (
	"context"

	"github.com/yoorita/currency-converter/app/controllers"
	"github.com/go-masonry/mortar/interfaces/monitor"
	converter "github.com/yoorita/currency-converter/api"
	"github.com/yoorita/currency-converter/app/validations"

	"github.com/go-masonry/mortar/interfaces/log"
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

func (impl *currencyConverterServiceImpl) Convert(ctx context.Context, req *converter.ConvertRequest) (*converter.ConvertResponse, error) {
	return &converter.ConvertResponse{}, nil
}
