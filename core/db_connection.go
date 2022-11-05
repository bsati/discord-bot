package core

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func dbConnect(connectionString string) *sql.DB {
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
