package clients

const (
	dbFilename = "exchange.sqlite.db"
	codesFilename = "exchange.sqlite.codes"


	createCodesTable = `
	CREATE TABLE IF NOT EXISTS codes(
		numericcode TEXT PRIMARY KEY,
		alphabeticcode TEXT NOT NULL UNIQUE,
		currency TEXT NOT NULL
	)
	`

	// SELECT
	getCountCodes = "SELECT COUNT(*) FROM codes"
	getCodeNumericValue = "SELECT numericcode FROM codes WHERE alphabeticcode = ?"
	// INSERT
	insertToCodes = "INSERT INTO codes(numericcode, alphabeticcode, currency) VALUES (?,?,?)"
)