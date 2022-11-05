package core

import (
	"database/sql"
)

type Env struct {
	DB *sql.DB
}

func BuildEnv(config *Config) *Env {
	return &Env{
		DB: dbConnect(config.DbConnectionString),
	}
}
