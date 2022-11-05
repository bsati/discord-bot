package core

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	env       *Env
	dgSession *discordgo.Session
}

func NewBot(config_path *string) (*Bot, error) {
	cfg := loadConfig(config_path)
	env := BuildEnv(&cfg)

	dg, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		return nil, err
	}

	return &Bot{env: env, dgSession: dg}, nil
}

func (b *Bot) Run() error {
	err := b.dgSession.Open()
	if err != nil {
		return err
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	b.dgSession.Close()
	return nil
}
