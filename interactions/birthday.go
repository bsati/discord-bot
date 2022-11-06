package interactions

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type BirthdayInteractions struct {
}

func (domain *BirthdayInteractions) CreateInteractions(session *discordgo.Session) *map[string][]*InteractionInfo {
	minMonth := 1.0
	interactions := []*discordgo.ApplicationCommand{
		{
			Name:        "birthday",
			Description: "Add / remove your birthday.",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "date",
					Description: "Your birthday in \"dd.mm.yyyy\" format, e.g. 01.01.1970.",
					Type:        discordgo.ApplicationCommandOptionString,
				},
				{
					Name:        "remove",
					Description: "Remove your birthday to stop receiving server messages.",
					Type:        discordgo.ApplicationCommandOptionBoolean,
				},
			},
		},
		{
			Name:        "birthdays",
			Description: "Get details about birthdays",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "member",
					Description: "Get a specific member's birthday",
					Type:        discordgo.ApplicationCommandOptionUser,
				},
				{
					Name:        "next",
					Description: "Get a list of the <next> upcoming birthdays",
					Type:        discordgo.ApplicationCommandOptionInteger,
					MinValue:    &minMonth,
					MaxValue:    12.0,
				},
				{
					Name:        "month",
					Description: "Get a list of all birthdays for the given month",
					Type:        discordgo.ApplicationCommandOptionInteger,
					MinValue:    &minMonth,
					MaxValue:    12.0,
				},
			},
		},
	}

	registeredInteractions := make(map[string][]*InteractionInfo)

	for _, guild := range session.State.Guilds {
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
	return &registeredInteractions
}

func (domain *BirthdayInteractions) CreateHandlers() *map[string]*InteractionHandler {
	handlers := make(map[string]*InteractionHandler)

	return &handlers
}
