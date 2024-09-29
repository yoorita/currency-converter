package mortar

import (
	"github.com/yoorita/currency-converter/app/clients"
	"go.uber.org/fx"
)

func SQLiteFxOptions() fx.Option {
	return fx.Options(
		fx.Provide(clients.CreateSQLiteClient),
	)
}