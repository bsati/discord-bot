package interactions

import (
	"github.com/bwmarrin/discordgo"
)

type InteractionRegistry struct {
	handlers               map[string]InteractionHandler
	registeredInteractions map[string]InteractionInfo
}

func InitInteractionRegistry(session *discordgo.Session) *InteractionRegistry {
	registry := InteractionRegistry{
		make(map[string]InteractionHandler),
		make(map[string]InteractionInfo),
	}

	registry.registerDomain(&BirthdayInteractions{}, session)

	return &registry
}

func (registry *InteractionRegistry) registerDomain(domain InteractionDomain, session *discordgo.Session) {
	_ = domain.CreateInteractions(session)
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
