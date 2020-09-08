package repositories

import "database/sql"

type database struct {
	conn *sql.DB
}

var DB database