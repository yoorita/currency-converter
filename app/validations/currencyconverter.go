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
		impl.validateField(ctx, req.GetCurrencyFrom()),
		impl.validateField(ctx, req.GetCurrencyTo()),
		func() error {
			if amount:=req.GetAmountFrom(); amount < 0 {
				impl.deps.Logger.WithField("amount_from", amount).Error(ctx, "negative value is not supported")
				return status.Errorf(codes.InvalidArgument, errorMessage("amount_from"))
			}
			return nil
		}(),
	)

	if currentError := multierr.Errors(combinedErrors); len(currentError) > 0 {
		err = currentError[0]
	}
	return
}

func (impl *currencyConverterValidationsImpl) validateField(ctx context.Context, object interface{}) (err error) {
	objType := reflect.TypeOf(object)
	if reflect.DeepEqual(object, reflect.Zero(objType).Interface()) {
		impl.deps.Logger.WithField(objType.Elem().Name(), object).Error(ctx, "invalid empty value")
		err = status.Errorf(codes.InvalidArgument, errorMessage(reflect.TypeOf(object).Elem().Name()))
	}
	return
}

func errorMessage(param string) string {
	return fmt.Sprintf("fill in the correct value for the parameter %v", param)
}