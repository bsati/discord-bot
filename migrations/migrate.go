package migrations

import (
	"fmt"
	"log"
	"os"

	"github.com/bsati/discord-bot/core"
)

func Migrate(fromVersion, targetVersion int) {
	cfg := core.LoadConfig(nil)
	db := core.DBConnect(cfg.DbConnectionString)
	for i := fromVersion; i < targetVersion+1; i++ {
		sqlFile, err := os.ReadFile(fmt.Sprintf("./migrations/sql/%d.sql", i))
		if err != nil {
			log.Panicf("Error reading sql file for version %d - Err: %s", i, err.Error())
		}
		sql := string(sqlFile)
		_, err = db.Query(sql)
		if err != nil {
			log.Panicf("Error reading sql file for version %d - Err: %s", i, err.Error())
		}
	}
}
