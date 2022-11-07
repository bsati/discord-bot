package interactions

import (
	"time"

	"github.com/bsati/discord-bot/services"
	"github.com/bwmarrin/discordgo"
)

type birthdayInteractions struct {
}

func (domain *birthdayInteractions) GetInteractions(session *discordgo.Session) []*discordgo.ApplicationCommand {
	minMonth := 1.0
	return []*discordgo.ApplicationCommand{
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
}

func (domain *birthdayInteractions) CreateHandlers(serviceRegistry *services.ServiceRegistry) *map[string]interactionHandler {
	handlers := make(map[string]interactionHandler)
	handlers["birthday"] = handleBirthday(serviceRegistry)
	return &handlers
}

func handleBirthday(birthdayService services.BirthdayService) interactionHandler {
	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
		options := interactionOptionsToMap(interaction)
		if _, ok := options["remove"]; ok {
			return birthdayService.RemoveBirthday(interaction.Member.User.ID)
		}
		if opt, ok := options["date"]; ok {
			parsed, err := time.Parse("02.01.2006", opt.StringValue())
			if err != nil {
				return newInteractionError("Invalid date format")
			}
			_, err = birthdayService.AddBirthday(interaction.Member.User.ID, parsed)
			interactionPrivateMessageResponse(session, interaction, "Birthday registered! 🎉🎉")
		}
		birthdayService.AddBirthday(interaction.Member.User.ID, time.Now())
		return nil
	}
}
