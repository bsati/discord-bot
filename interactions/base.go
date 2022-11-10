package interactions

import (
	"log"

	"github.com/bsati/discord-bot/daos"
	"github.com/bwmarrin/discordgo"
)

type interactionRegistry struct {
	handlers               map[string]interactionHandler
	registeredInteractions map[string][]*interactionInfo
}

// InitInteractionHandling initializes the interaction handling by constructing
// an interactionRegistry and adding Discord handlers for events
func InitInteractionHandling(session *discordgo.Session, dao *daos.DAO) {
	registry := interactionRegistry{
		handlers: make(map[string]interactionHandler),
	}

	domains := []interactionDomain{&birthdayInteractions{}, &generalInteractions{}}

	interactions := []*discordgo.ApplicationCommand{}
	for _, domain := range domains {
		domainInteractions := registry.registerDomain(domain, session, dao)
		interactions = append(interactions, domainInteractions...)
	}

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := registry.handlers[i.ApplicationCommandData().Name]; ok {
			err := handler(s, i)
			if err != nil {
				interactionPrivateMessageResponse(s, i, "Error", err.Error())
			}
		}
	})

	session.AddHandler(func(s *discordgo.Session, e *discordgo.Ready) {
		registeredInteractions := make(map[string][]*interactionInfo, len(interactions))

		for _, guild := range session.State.Guilds {
			log.Printf("Initializing Interactions for Guild with ID: %s, Name: %s\n", guild.ID, guild.Name)
			for _, domain := range domains {
				domain.InitGuild(s, guild, dao)
			}
			registeredInteractions[guild.ID] = make([]*interactionInfo, len(interactions))
			for i, v := range interactions {
				cmd, err := session.ApplicationCommandCreate(session.State.User.ID, guild.ID, v)
				if err != nil {
					log.Panicf("Cannot create '%v' command: %v", v.Name, err)
				}
				registeredInteractions[guild.ID][i] = &interactionInfo{
					AppID:   cmd.ApplicationID,
					GuildID: cmd.GuildID,
					CmdID:   cmd.ID,
				}
			}
		}

		registry.registeredInteractions = registeredInteractions
	})
}

func (registry *interactionRegistry) registerDomain(domain interactionDomain, session *discordgo.Session, dao *daos.DAO) []*discordgo.ApplicationCommand {
	interactions := domain.GetInteractions(session)

	handlers := domain.CreateHandlers(dao)
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

type interactionInfo struct {
	AppID   string
	GuildID string
	CmdID   string
}

type interactionDomain interface {
	// GetInteractions returns a list of possible interactions for the domain that can be contructed
	// for different guilds
	GetInteractions(session *discordgo.Session) []*discordgo.ApplicationCommand
	// CreateHandlers returns a lsit of handlers for the interactions supplied by GetInteractions
	CreateHandlers(dao *daos.DAO) *map[string]interactionHandler
	// InitGuild is triggered for every guild in the ReadyEvent, since some domains e.g. birthdays
	// need to initialize timers etc.
	InitGuild(session *discordgo.Session, guild *discordgo.Guild, dao *daos.DAO)
}

type interactionHandler func(session *discordgo.Session, interaction *discordgo.InteractionCreate) error

type interactionBaseError struct {
	message string
}

func (e *interactionBaseError) Error() string {
	return e.message
}

func newInteractionError(message string) error {
	return &interactionBaseError{
		message: message,
	}
}
