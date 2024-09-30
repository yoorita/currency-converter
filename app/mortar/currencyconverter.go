package mortar

import (
	"context"

	serverInt "github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/go-masonry/mortar/providers/groups"
	converter "github.com/yoorita/currency-converter/api"
	"github.com/yoorita/currency-converter/app/controllers"
	"github.com/yoorita/currency-converter/app/data"
	"github.com/yoorita/currency-converter/app/services"
	"github.com/yoorita/currency-converter/app/validations"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type converterServiceDeps struct {
	fx.In

	CurrencyConverter converter.CurrencyConverterServer
}

func ServiceAPIsAndOtherDependenciesFxOption() fx.Option {
	return fx.Options(
		// GRPC Service APIs registration
		fx.Provide(fx.Annotated{
			Group:  groups.GRPCServerAPIs,
			Target: serviceGRPCServiceAPIs,
		}),
		// GRPC Gateway Generated Handlers registration
		fx.Provide(fx.Annotated{
			Group:  groups.GRPCGatewayGeneratedHandlers + ",flatten", // "flatten" does this [][]serverInt.GRPCGatewayGeneratedHandlers -> []serverInt.GRPCGatewayGeneratedHandlers
			Target: serviceGRPCGatewayHandlers,
		}),
		// All other dependencies
		serviceDependencies(),
	)
}

func serviceGRPCServiceAPIs(deps converterServiceDeps) serverInt.GRPCServerAPI {
	return func(srv *grpc.Server) {
		converter.RegisterCurrencyConverterServer(srv, deps.CurrencyConverter)
		// Any additional gRPC Implementations should be called here
	}
}

func serviceGRPCGatewayHandlers() []serverInt.GRPCGatewayGeneratedHandlers {
	return []serverInt.GRPCGatewayGeneratedHandlers{
		// Register service REST API
		func(mux *runtime.ServeMux, localhostEndpoint string) error {
			return converter.RegisterCurrencyConverterHandlerFromEndpoint(
				context.Background(), 
				mux, 
				localhostEndpoint, 
				[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
			)
		},
		// Any additional gRPC gateway registrations should be called here
	}
}

func serviceDependencies() fx.Option {
	return fx.Provide(
		services.CreateCurrencyConverterdService,
		controllers.CreateCurrencyConverterController,
		validations.CreateCurrencyConverterValidations,
		data.CreateCurrencyConverterDao,
	)
}
