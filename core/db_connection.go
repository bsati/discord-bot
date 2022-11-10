package core

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// DBConnect creates a *sql.DB instance connected to the postgres database
// with the given connectionString or panics if there are any errors
func DBConnect(connectionString string) *sql.DB {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
	if err = db.Ping(); err != nil {
		log.Println(err.Error())
		panic(err)
	}
	return db
}
