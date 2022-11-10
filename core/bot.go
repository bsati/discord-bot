package core

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bsati/discord-bot/daos"
	"github.com/bsati/discord-bot/interactions"
	"github.com/bwmarrin/discordgo"
)

// Bot serves as the main struct to hold relevant variables
type Bot struct {
	env       *Env
	dgSession *discordgo.Session
}

// NewBot creates a Bot instance by loading a config with the specified path
// and initializes the database connection, interaction handling, ...
func NewBot(config_path *string) (*Bot, error) {
	cfg := LoadConfig(config_path)
	env := BuildEnv(&cfg)

	dg, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		return nil, err
	}
	log.Println("Bot connected")
	log.Println("Initializing Services")
	dao := daos.NewDAO(env.DB)
	log.Println("Services initialized")
	log.Println("Initializing InteractionRegistry")
	interactions.InitInteractionHandling(dg, dao)
	log.Println("InteractionRegistry initiliazed")

	return &Bot{env: env, dgSession: dg}, nil
}

// Run starts the Discord session and creates a channel to gracefully close
// the connection on Ctrl+C
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
