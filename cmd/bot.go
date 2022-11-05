package main

import (
	"github.com/bsati/discord-bot/core"
)

func main() {
	bot, err := core.NewBot(nil)
	if err != nil {
		panic(err)
	}

	err = bot.Run()
	if err != nil {
		panic(err)
	}
}
