package migrations

import (
	"fmt"
	"log"
	"os"

	"github.com/bsati/discord-bot/core"
)

// Migrate serves as a small migration scripting engine.
// SQLs are read from the folder "./migrations/sql" and are applied in order
// by their version from <fromVersion> to <targetVersion>.
func Migrate(fromVersion, targetVersion int) {
	cfg := core.LoadConfig(nil)
	dbName := os.Getenv("DB_NAME")
	if fromVersion == 0 && targetVersion >= 0 {
		db := core.DBConnect(cfg.DbConnectionString)
		_, err := db.Exec("CREATE DATABASE " + dbName)
		if err != nil {
			log.Panicf("Error creating database - Err: %s\n", err.Error())
		}
		db.Close()
		fromVersion++
	}

	db := core.DBConnect(fmt.Sprintf("%s dbname=%s", cfg.DbConnectionString, dbName))

	for i := fromVersion; i < targetVersion+1; i++ {
		sqlFile, err := os.ReadFile(fmt.Sprintf("./migrations/sql/%d.sql", i))
		if err != nil {
			log.Panicf("Error reading sql file for version %d - Err: %s\n", i, err.Error())
		}
		sql := string(sqlFile)
		_, err = db.Query(sql)
		if err != nil {
			log.Panicf("Error running sql file for version %d - Err: %s\n", i, err.Error())
		}
	}

	db.Close()
}
