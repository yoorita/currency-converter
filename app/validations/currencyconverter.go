package validations

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-masonry/mortar/interfaces/log"
	converter "github.com/yoorita/currency-converter/api"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (impl *currencyConverterValidationsImpl) ValidateCurruncyConvertRequest(ctx context.Context, req *converter.ConvertRequest) (err error) {
	combinedErrors := multierr.Combine(
		validateField(req.CurrencyFrom),
		validateField(req.CurrencyTo),
		validateField(req.AmountFrom))

	if currentError := multierr.Errors(combinedErrors); len(currentError) > 0 {
		err = currentError[0]
	}
	return
}

func validateField(object interface{}) (err error) {
	if isEmpty(object) {
		err = status.Errorf(codes.InvalidArgument, fmt.Sprintf("Input parameter %v cannot be empty", reflect.TypeOf(object).Elem().Name()))
	}
	return
}

func isEmpty(object interface{}) bool {
	if object == nil {
		return true
	}

	zero := reflect.Zero(reflect.TypeOf(object))
	return reflect.DeepEqual(object, zero.Interface())
}
