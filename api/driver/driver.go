package driver

import (
	"database/sql"
	"fmt"
	"log"
)

func Connect(dbDriver, dbUrl string) *sql.DB {
	db, err := sql.Open(dbDriver, dbUrl)
	if err != nil {
		fmt.Printf("\nCannot connect to %s database", dbDriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("\nWe are connected to the %s database", dbDriver)
	}
	return db
}
