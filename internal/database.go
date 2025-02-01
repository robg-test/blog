package internal

import (
	"database/sql"
	"log"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("libsql", dbPath)
	log.Print("database initialized")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
		return err
	}
	return nil
}
