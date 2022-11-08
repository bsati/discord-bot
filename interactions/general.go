package interactions

import (
	"log"

	"github.com/bsati/discord-bot/daos"
	"github.com/bwmarrin/discordgo"
)

type generalInteractions struct {
}

func (domain *generalInteractions) GetInteractions(session *discordgo.Session) []*discordgo.ApplicationCommand {
	manageServer := int64(discordgo.PermissionManageServer)
	return []*discordgo.ApplicationCommand{
		{
			Name:                     "designate_channel",
			Description:              "Designate the channel to be a bot channel.",
			Type:                     discordgo.ChatApplicationCommand,
			DefaultMemberPermissions: &manageServer,
		},
	}
}

func (domain *generalInteractions) CreateHandlers(dao *daos.DAO) *map[string]interactionHandler {
	handlers := make(map[string]interactionHandler)
	handlers["designate_channel"] = handleChannelDesignation(dao)
	return &handlers
}

func handleChannelDesignation(dao *daos.DAO) interactionHandler {
	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
		err := dao.SetBotChannelForGuild(interaction.ChannelID, interaction.GuildID)
		if err != nil {
			log.Printf("Error setting designated bot channel %s for guild with id %s: %v\n", interaction.ChannelID, interaction.GuildID, err)
			return newInteractionError("Error setting designated channel.")
		}
		interactionPrivateMessageResponse(session, interaction, "Bot channel successfully set.")
		return nil
	}
}

func (domain *generalInteractions) InitGuild(session *discordgo.Session, guild *discordgo.Guild, dao *daos.DAO) {
}
