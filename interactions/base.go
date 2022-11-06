package interactions

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type InteractionRegistry struct {
	handlers               map[string]*InteractionHandler
	registeredInteractions map[string][]*InteractionInfo
}

func InitInteractionRegistry(session *discordgo.Session) *InteractionRegistry {
	registry := InteractionRegistry{
		make(map[string]*InteractionHandler),
		make(map[string][]*InteractionInfo),
	}

	registry.registerDomain(&BirthdayInteractions{}, session)

	return &registry
}

func (registry *InteractionRegistry) registerDomain(domain InteractionDomain, session *discordgo.Session) {
	interactions := domain.CreateInteractions(session)
	for key, val := range *interactions {
		if current, ok := registry.registeredInteractions[key]; ok {
			registry.registeredInteractions[key] = append(current, val...)
		} else {
			registry.registeredInteractions[key] = val
		}
	}
	handlers := domain.CreateHandlers()
	for key, val := range *handlers {
		if _, ok := registry.handlers[key]; ok {
			log.Printf("Interaction handler \"%s\" has already been registered and is about to be reregistered by domain \"%T\", skipping.\n", key, domain)
		} else {
			registry.handlers[key] = val
		}
	}
	log.Printf("Registered domain \"%T\"\n", domain)
}

type InteractionInfo struct {
	AppID   string
	GuildID string
	CmdID   string
}

type InteractionDomain interface {
	CreateInteractions(session *discordgo.Session) *map[string][]*InteractionInfo
	CreateHandlers() *map[string]*InteractionHandler
}

type InteractionHandler func(session *discordgo.Session, interaction *discordgo.InteractionCreate)
