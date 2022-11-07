package main

import (
	"log"
	"os"
	"strconv"

	"github.com/bsati/discord-bot/migrations"
)

func main() {
	targetVersion, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Panicf("Invalid target_version specified: %s", err.Error())
	}
	fromVersion := 0
	if len(os.Args) > 2 {
		fromVersion, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Panicf("Invalid from_version specified: %s", err.Error())
		}
	}

	migrations.Migrate(fromVersion, targetVersion)
}
