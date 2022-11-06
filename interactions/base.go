package interactions

import (
	"log"

	"github.com/bsati/discord-bot/services"
	"github.com/bwmarrin/discordgo"
)

type InteractionRegistry struct {
	handlers               map[string]InteractionHandler
	registeredInteractions map[string][]*InteractionInfo
}

func InitInteractionHandling(session *discordgo.Session, serviceRegistry *services.ServiceRegistry) {
	registry := InteractionRegistry{
		make(map[string]InteractionHandler),
		make(map[string][]*InteractionInfo),
	}

	birthdayInteractions := registry.registerDomain(&BirthdayInteractions{}, session, serviceRegistry)

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := registry.handlers[i.ApplicationCommandData().Name]; ok {
			err := handler(s, i)
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   discordgo.MessageFlagsEphemeral,
						Content: err.Error(),
					},
				})
			}
		}
	})

	interactions := make([]*discordgo.ApplicationCommand, 0)
	interactions = append(interactions, birthdayInteractions...)

	session.AddHandler(func(s *discordgo.Session, e *discordgo.Ready) {
		registeredInteractions := make(map[string][]*InteractionInfo, len(interactions))

		for _, guild := range session.State.Guilds {
			log.Printf("Initializing Interactions for Guild with ID: %s, Name: %s\n", guild.ID, guild.Name)
			registeredInteractions[guild.ID] = make([]*InteractionInfo, len(interactions))
			for i, v := range interactions {
				cmd, err := session.ApplicationCommandCreate(session.State.User.ID, guild.ID, v)
				if err != nil {
					log.Panicf("Cannot create '%v' command: %v", v.Name, err)
				}
				registeredInteractions[guild.ID][i] = &InteractionInfo{
					AppID:   cmd.ApplicationID,
					GuildID: cmd.GuildID,
					CmdID:   cmd.ID,
				}
			}
		}
	})
}

func (registry *InteractionRegistry) registerDomain(domain InteractionDomain, session *discordgo.Session, serviceRegistry *services.ServiceRegistry) []*discordgo.ApplicationCommand {
	interactions := domain.GetInteractions(session)

	handlers := domain.CreateHandlers(serviceRegistry)
	for key, val := range *handlers {
		if _, ok := registry.handlers[key]; ok {
			log.Printf("Interaction handler \"%s\" has already been registered and is about to be reregistered by domain \"%T\", skipping.\n", key, domain)
		} else {
			registry.handlers[key] = val
		}
	}
	log.Printf("Registered domain \"%T\"\n", domain)
	return interactions
}

type InteractionInfo struct {
	AppID   string
	GuildID string
	CmdID   string
}

type InteractionDomain interface {
	GetInteractions(session *discordgo.Session) []*discordgo.ApplicationCommand
	CreateHandlers(serviceRegistry *services.ServiceRegistry) *map[string]InteractionHandler
}

type InteractionHandler func(session *discordgo.Session, interaction *discordgo.InteractionCreate) error
