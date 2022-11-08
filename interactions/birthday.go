package interactions

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bsati/discord-bot/daos"
	"github.com/bsati/discord-bot/models"
	"github.com/bwmarrin/discordgo"
)

type birthdayInteractions struct {
}

func (domain *birthdayInteractions) GetInteractions(session *discordgo.Session) []*discordgo.ApplicationCommand {
	minMonth := 1.0
	return []*discordgo.ApplicationCommand{
		{
			Name:        "birthday",
			Description: "Add your birthday so everybody can congratulate you!",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "date",
					Description: "Your birthday in \"dd.mm.yyyy\" format, e.g. 01.01.1970.",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:        "birthday_remove",
			Description: "Remove your registered birthday to stop receiving messages.",
			Type:        discordgo.ChatApplicationCommand,
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
					Description: "Get a list of upcoming birthdays for the <next> months",
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

func (domain *birthdayInteractions) CreateHandlers(dao *daos.DAO) *map[string]interactionHandler {
	handlers := make(map[string]interactionHandler)
	handlers["birthday"] = handleBirthday(dao)
	handlers["birthday_remove"] = handleRemoveBirthday(dao)
	handlers["birthdays"] = handleBirthdayList(dao)
	return &handlers
}

func (domain *birthdayInteractions) InitGuild(session *discordgo.Session, guild *discordgo.Guild, dao *daos.DAO) {
	handleBirthdaysOfDay(session, guild, dao)

	now := time.Now()
	year, month, day := now.Date()
	nextMorning := time.Date(year, month, day+1, 8, 0, 0, 0, now.Location())

	go func() {
		time.Sleep(nextMorning.Sub(now))
		handleBirthdaysOfDay(session, guild, dao)

		ticker := time.NewTicker(24 * time.Hour)

		go func() {
			for range ticker.C {
				handleBirthdaysOfDay(session, guild, dao)
			}
		}()
	}()
}

func handleBirthdaysOfDay(session *discordgo.Session, guild *discordgo.Guild, dao *daos.DAO) {
	birthdays, err := dao.GetBirthdaysByDay(guild.ID, time.Now())
	if err != nil {
		log.Printf("Error initializing birthdays for guild id %s: %v\n", guild.ID, err)
	}

	if len(birthdays) == 0 {
		return
	}

	usernames := make([]string, len(birthdays))
	for i, birthday := range birthdays {
		user, err := session.User(birthday.UserId)
		if err != nil {
			log.Printf("Error fetching user with id %s: %v\n", birthday.UserId, err)
			continue
		}
		usernames[i], err = getUsername(session, guild.ID, user)
		if err != nil {
			log.Printf("Error retrieving username for user with id %s: %v\n", birthday.UserId, err)
			continue
		}
	}

	channels, err := dao.GetBotChannelByGuild(guild.ID)
	var channelId string
	if err != nil {
		log.Printf("Error fetching designated channels for guild with %s: %v. Using base channel.\n", guild.ID, err)
		channelId = guild.Channels[0].ID
	} else if len(channels) == 0 {
		log.Printf("No designated channels found for guild with id %s. Using base channel.\n", guild.ID)
		channelId = (*(guild.Channels[0])).ID
	} else {
		channelId = channels[0]
	}

	var usernamesFormatted string
	usernameCount := len(usernames)
	if usernameCount == 1 {
		usernamesFormatted = usernames[0]
	} else if usernameCount == 2 {
		usernamesFormatted = fmt.Sprintf("%s and %s", usernames[0], usernames[1])
	} else {
		usernamesFormatted = fmt.Sprintf("%s and %s", usernames[0:(usernameCount-2)], usernames[usernameCount-1])
	}

	session.ChannelMessageSend(channelId, fmt.Sprintf("Happy birthday to: %s! 🎉🎉", usernamesFormatted))
}

func handleBirthday(dao *daos.DAO) interactionHandler {
	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
		parsed, err := time.Parse("02.01.2006", interaction.ApplicationCommandData().Options[0].StringValue())
		if err != nil {
			return newInteractionError("The date format you entered is invalid.")
		}
		_, err = dao.AddBirthday(interaction.Member.User.ID, interaction.GuildID, parsed)
		if err != nil {
			if strings.HasPrefix(err.Error(), "pq: duplicate key") {
				return newInteractionError("Your birthday has already been added.")
			}
			log.Printf("Error adding birthday: %v\n", err)
			return newInteractionError("Unknown error adding your birthday.")
		}
		interactionPrivateMessageResponse(session, interaction, "Birthday registered! 🎉🎉")
		return nil
	}
}

func handleRemoveBirthday(dao *daos.DAO) interactionHandler {
	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
		err := dao.RemoveBirthday(interaction.Member.User.ID, interaction.GuildID)
		if err != nil {
			log.Printf("Error removing birthday: %v\n", err)
			return newInteractionError("Error removing your birthday.")
		}
		interactionPrivateMessageResponse(session, interaction, "Your birthday has been removed!")
		return nil
	}
}

func handleBirthdayList(dao *daos.DAO) interactionHandler {
	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate) error {
		optionMap := interactionOptionsToMap(interaction)
		if option, ok := optionMap["member"]; ok {
			user := option.UserValue(session)
			birthday, err := dao.GetBirthdayByUserId(user.ID, interaction.GuildID)
			if err != nil {
				return newInteractionError("The user has not registered his birthday.")
			}
			username, err := getUsername(session, interaction.GuildID, user)
			if err != nil {
				return err
			}
			today := time.Now()
			if birthday.Date.Day() == today.Day() && birthday.Date.Month() == today.Month() {
				interactionMessageResponse(session, interaction, fmt.Sprintf("%s's birthday is today! 🎉🎉", username))
				return nil
			}
			interactionMessageResponse(session, interaction, fmt.Sprintf("%s's birthday is on %s!", username, birthday.Date.Format("02.01.2006")))
			return nil
		}
		if option, ok := optionMap["next"]; ok {
			nextMonths := int(option.IntValue())
			birthdays, err := dao.GetUpcomingBirthdaysForMonths(interaction.GuildID, nextMonths, time.Now())
			if err != nil {
				log.Printf("Error retrieving birthdays for upcoming months: %v\n", err)
				return newInteractionError("Unknown error occured.")
			}
			interactionMessageResponse(session, interaction, formatBirthdaysToMessage(session, interaction.GuildID, birthdays))
			return nil
		}
		if option, ok := optionMap["month"]; ok {
			birthdays, err := dao.GetBirthdaysByMonth(interaction.GuildID, int(option.IntValue()))
			if err != nil {
				log.Printf("Error retrieving birthdays for month %d: %v\n", option.IntValue(), err)
			}
			interactionMessageResponse(session, interaction, formatBirthdaysToMessage(session, interaction.GuildID, birthdays))
		}
		interactionMessageResponse(session, interaction, "Please select one option")
		return nil
	}
}

func formatBirthdaysToMessage(s *discordgo.Session, guildId string, birthdays []models.Birthday) string {
	var builder strings.Builder
	for _, birthday := range birthdays {
		user, err := s.User(birthday.UserId)
		if err != nil {
			log.Printf("Error fetching user with id %s: %v", birthday.UserId, err)
			continue
		}
		username, err := getUsername(s, guildId, user)
		if err != nil {
			log.Printf("Error getting username for user with id %s and guild %s", user.ID, guildId)
			continue
		}
		builder.WriteString(birthday.Date.Format("02.01.2006"))
		builder.WriteString(": ")
		builder.WriteString(username)
		builder.WriteRune('\n')
	}
	return builder.String()
}
