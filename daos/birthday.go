package daos

import (
	"database/sql"
	"time"

	"github.com/bsati/discord-bot/models"
)

type BirthdayDAO interface {
	GetBirthdays(guildId string, limit int, from time.Time) ([]models.Birthday, error)
	GetBirthdayByUserId(userId, guildId string) (models.Birthday, error)
	GetUpcomingBirthdaysForMonths(guildId string, months int, from time.Time) ([]models.Birthday, error)
	GetBirthdaysByMonth(guildId string, month int) ([]models.Birthday, error)
	GetBirthdaysByDay(guildId string, day time.Time) ([]models.Birthday, error)
	AddBirthday(userId, guildId string, date time.Time) (*models.Birthday, error)
	RemoveBirthday(userId, guildId string) error
}

type birthdayDAOSQL struct {
	db *sql.DB
}

func NewBirthdayDAO(db *sql.DB) BirthdayDAO {
	return &birthdayDAOSQL{db: db}
}

func (br *birthdayDAOSQL) GetBirthdays(guildId string, limit int, from time.Time) ([]models.Birthday, error) {
	rows, err := br.db.Query(`SELECT b.id, b.user_id, b.date FROM birthdays AS b JOIN user_guild ON b.user_id = user_guild.user_id
	WHERE guild_id = $3 AND ((EXTRACT(MONTH FROM date) > $1) 
	OR (EXTRACT(MONTH FROM date) == $1 AND EXTRACT(DAY FROM date) >= $2)) ORDER BY EXTRACT(MONTH FROM date) ASC, EXTRACT(DAY FROM date) ASC`, from.Month(), from.Day(), guildId)
	return scanBirthdayRows(rows, err)
}

func (br *birthdayDAOSQL) GetUpcomingBirthdaysForMonths(guildId string, months int, from time.Time) ([]models.Birthday, error) {
	var rows *sql.Rows
	var err error
	fromMonth := int(from.Month())
	if fromMonth+months > 12 {
		rows, err = br.db.Query(`SELECT b.id, b.user_id, b.date FROM birthdays AS b JOIN user_guild ON b.user_id = user_guild.user_id
		WHERE guild_id = $1 AND ((EXTRACT(MONTH FROM date) > $2) OR (EXTRACT(MONTH FROM date) <= $3)) ORDER BY EXTRACT(MONTH FROM date) ASC, EXTRACT(DAY FROM date) ASC`, guildId, fromMonth, fromMonth+months-12)
	} else {
		rows, err = br.db.Query(`SELECT b.id, b.user_id, b.date FROM birthdays AS b JOIN user_guild ON b.user_id = user_guild.user_id
		WHERE guild_id = $1 AND (EXTRACT(MONTH FROM date) > $2) AND (EXTRACT(MONTH FROM date) <= $3) ORDER BY EXTRACT(MONTH FROM date) ASC, EXTRACT(DAY FROM date) ASC`, guildId, fromMonth, fromMonth+months)
	}
	return scanBirthdayRows(rows, err)
}

func (br *birthdayDAOSQL) GetBirthdaysByMonth(guildId string, month int) ([]models.Birthday, error) {
	rows, err := br.db.Query(`SELECT b.id, b.user_id, b.date FROM birthdays AS b JOIN user_guild ON b.user_id = user_guild.user_id
	WHERE guild_id = $1 AND EXTRACT(MONTH FROM date) = $2 ORDER BY EXTRACT(MONTH FROM date) ASC, EXTRACT(DAY FROM date) ASC`, guildId, month)
	return scanBirthdayRows(rows, err)
}

func (br *birthdayDAOSQL) GetBirthdaysByDay(guildId string, day time.Time) ([]models.Birthday, error) {
	rows, err := br.db.Query(`SELECT b.id, b.user_id, b.date FROM birthdays AS b JOIN user_guild ON b.user_id = user_guild.user_id
	WHERE guild_id = $1 AND EXTRACT(MONTH from date) = $2 AND EXTRACT(DAY from date) = $3`, guildId, day.Month(), day.Day())
	return scanBirthdayRows(rows, err)
}

func scanBirthdayRows(rows *sql.Rows, err error) ([]models.Birthday, error) {
	var result []models.Birthday
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var birthday models.Birthday
		err = rows.Scan(&birthday.Id, &birthday.UserId, &birthday.Date)
		if err == nil {
			result = append(result, birthday)
		}
	}
	return result, nil
}

func (br *birthdayDAOSQL) GetBirthdayByUserId(userId, guildId string) (models.Birthday, error) {
	var result models.Birthday
	err := br.db.QueryRow(`SELECT b.id, b.user_id, b.date FROM birthdays AS b JOIN user_guild ON b.user_id = user_guild.user_id
	WHERE guild_id = $1 AND b.user_id = $2`, guildId, userId).Scan(&result.Id, &result.UserId, &result.Date)
	return result, err
}

func (br *birthdayDAOSQL) AddBirthday(userId, guildId string, date time.Time) (*models.Birthday, error) {
	result := &models.Birthday{
		UserId: userId,
		Date:   date,
	}

	tx, err := br.db.Begin()
	if err != nil {
		return result, err
	}

	defer tx.Rollback()

	err = tx.QueryRow(`INSERT INTO birthdays (user_id, date) VALUES ($1, $2) RETURNING Id`, userId, date).Scan(&result.Id)
	if err != nil {
		return result, err
	}
	_, err = tx.Exec(`INSERT INTO user_guild (user_id, guild_id) VALUES ($1, $2)`, userId, guildId)
	if err != nil {
		return result, err
	}
	return result, tx.Commit()
}

func (br *birthdayDAOSQL) RemoveBirthday(userId, guildId string) error {
	tx, err := br.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM birthdays WHERE user_id = $1`, userId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM user_guild WHERE user_id = $1 AND guild_id = $2`, userId, guildId)
	if err != nil {
		return err
	}

	return tx.Commit()
}
