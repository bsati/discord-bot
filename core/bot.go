package core

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bsati/discord-bot/interactions"
	"github.com/bsati/discord-bot/services"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	env       *Env
	dgSession *discordgo.Session
}

func NewBot(config_path *string) (*Bot, error) {
	cfg := LoadConfig(config_path)
	env := BuildEnv(&cfg)

	dg, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		return nil, err
	}
	log.Println("Bot connected")
	log.Println("Initializing Services")
	serviceRegistry := services.InitServices(env.DB)
	log.Println("Services initialized")
	log.Println("Initializing InteractionRegistry")
	interactions.InitInteractionHandling(dg, serviceRegistry)
	log.Println("InteractionRegistry initiliazed")

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
