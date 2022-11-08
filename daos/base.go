package daos

import (
	"database/sql"
)

type DAO struct {
	db *sql.DB
	BirthdayDAO
}

func NewDAO(db *sql.DB) *DAO {
	return &DAO{
		db:          db,
		BirthdayDAO: &birthdayDAOSQL{db},
	}
}

func (dao *DAO) GetBotChannelByGuild(guildId string) ([]string, error) {
	rows, err := dao.db.Query(`SELECT channel_id FROM guild_bot_channel WHERE guild_id = $1`, guildId)
	var result []string
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var channelId string
		err = rows.Scan(&channelId)
		if err == nil {
			result = append(result, channelId)
		}
	}
	return result, nil
}

func (dao *DAO) SetBotChannelForGuild(channelId, guildId string) error {
	var foundChannelId string
	err := dao.db.QueryRow(`SELECT channel_id FROM guild_bot_channel WHERE channel_id = $1 AND guild_id = $2`, channelId, guildId).Scan(&foundChannelId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			_, err = dao.db.Exec(`INSERT INTO guild_bot_channel (channel_id, guild_id) VALUES ($1, $2)`, channelId, guildId)
		}
		return err
	}
	_, err = dao.db.Exec(`UPDATE guild_bot_channel SET channel_id = $1 WHERE guild_id = $2 AND channel_id = $3`, channelId, guildId, foundChannelId)
	return err
}
