package common

import "database/sql"

type Mission struct {
	ID       int64
	CatID    sql.NullInt64 // Changed to sql.NullInt64 to handle NULL values
	Complete bool
	Targets  []Target
}
