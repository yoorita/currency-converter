package clients

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"

	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yoorita/currency-converter/model"
	"go.uber.org/fx"
)

type (
	LazySQLCLient struct {
		Client SQLiteClient
	}
	SQLiteClient interface {
		GetCurrencyCode(context.Context, string) (string, error)
	}

	sqliteClientImplDeps struct {
		fx.In

		Logger    log.Logger
		Lifecycle fx.Lifecycle
		Config    cfg.Config
	}

	sqliteClientImpl struct {
		deps sqliteClientImplDeps
		db   *sql.DB
	}
)

func CreateSQLiteClient(deps sqliteClientImplDeps) (client *LazySQLCLient, err error) {
	var sqldb *sql.DB
	var clientPtr = new(LazySQLCLient)
	dbFileString := deps.Config.Get(dbFilename).String()

	deps.Lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) (err error) {
				if sqldb, err = sql.Open("sqlite3", dbFileString); err != nil {
					deps.Logger.WithError(err).Error(ctx, "could not connect to database")
					return
				}

				if _, err = sqldb.Exec(createCodesTable); err != nil {
					deps.Logger.WithError(err).Error(ctx, "couldn't create codes table")
					return
				}

				clientPtr.Client = &sqliteClientImpl{
					deps: deps,
					db: sqldb,
				}

				deps.Logger.Info(ctx, "connected to DB")

				row := sqldb.QueryRow(getCountCodes)
				var counted float32
				err = row.Scan(&counted)
			
				if err != nil {
					deps.Logger.WithError(err).Error(ctx, "couldn't get codes row amount from table")
					return
				}

				if counted > 0 {
					return
				}

				deps.Logger.Info(ctx, "codes table is empty - fetching from json")

				codesByte, err := os.ReadFile(deps.Config.Get(codesFilename).String())

				if err != nil {
					deps.Logger.WithError(err).Error(ctx, "couldn't open json with codes")
					return
				}
			
				var codes []model.CurrencyCodeModel
			
				if err = json.Unmarshal(codesByte, &codes); err != nil {
					deps.Logger.WithError(err).Error(ctx, "couldn't parse json with codes")
					return
				}

				statement, err := sqldb.Prepare(insertToCodes)

				if err != nil {
					deps.Logger.WithError(err).Error(ctx, "couldn't prepare for sql query")
					return
				}
			
				defer statement.Close()

				for _, code := range codes {
					_, err = statement.Exec(code.NumericCode, code.AlphabeticCode, code.Currency)
				}

				return
			},
			OnStop: func(ctx context.Context) error {
				return sqldb.Close()
			},
	})

	client = clientPtr
	return 
}

func (impl *sqliteClientImpl) GetCurrencyCode(ctx context.Context, alphabeticcode string) (result string, err error) {
	row := impl.db.QueryRow(getCodeNumericValue, alphabeticcode)
	err = row.Scan(&result)

	if err != nil {
		impl.deps.Logger.WithError(err).Error(ctx, "couldn't fetch numeric code by alphabetic code")
		return
	}

	impl.deps.Logger.WithField("from code", result).Info(ctx, "got code")
	return
}