package interactions

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func interactionOptionsToMap(interaction *discordgo.InteractionCreate) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, option := range options {
		optionMap[option.Name] = option
	}
	return optionMap
}

func interactionMessageResponse(s *discordgo.Session, i *discordgo.InteractionCreate, title, message string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       title,
					Description: message,
					Color:       16705372,
				},
			},
		},
	})
}

func interactionPrivateMessageResponse(s *discordgo.Session, i *discordgo.InteractionCreate, title, message string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       title,
					Description: message,
					Color:       16705372,
				},
			},
		},
	})
}

func getUsername(s *discordgo.Session, guildId string, user *discordgo.User) (string, error) {
	var username string
	member, err := s.GuildMember(guildId, user.ID)
	if err != nil {
		log.Printf("Error retrieving guild for the interaction: %v", err)
		return username, newInteractionError("Unknown error occured.")
	}
	username = member.Nick
	if username == "" {
		username = user.Username
	}
	return username, nil
}
