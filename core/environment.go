package core

import (
	"database/sql"
)

// Env represents a global environment that is needed
// in different parts of the app
type Env struct {
	DB *sql.DB
}

// BuildEnv creates an *Env instance with the given config
func BuildEnv(config *Config) *Env {
	return &Env{
		DB: DBConnect(config.DbConnectionString),
	}
}
