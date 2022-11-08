package daos

import (
	"database/sql"
	"time"

	"github.com/bsati/discord-bot/models"
)

type BirthdayDAO interface {
	GetBirthdays(limit int, from time.Time) ([]models.Birthday, error)
	GetBirthdayByUserId(userId string) (models.Birthday, error)
	GetUpcomingBirthdaysForMonths(months int, from time.Time) ([]models.Birthday, error)
	GetBirthdaysOfMonth(month int) ([]models.Birthday, error)
	AddBirthday(userId string, date time.Time) (*models.Birthday, error)
	RemoveBirthday(userId string) error
}

type birthdayDAOSQL struct {
	db *sql.DB
}

func NewBirthdayDAO(db *sql.DB) BirthdayDAO {
	return &birthdayDAOSQL{db: db}
}

func (br *birthdayDAOSQL) GetBirthdays(limit int, from time.Time) ([]models.Birthday, error) {
	rows, err := br.db.Query(`SELECT * FROM birthdays WHERE (EXTRACT(MONTH FROM date) > $1) 
	OR (EXTRACT(MONTH FROM date) == $1 AND EXTRACT(DAY FROM date) >= $2)`, from.Month(), from.Day())
	return scanBirthdayRows(rows, err)
}

func (br *birthdayDAOSQL) GetUpcomingBirthdaysForMonths(months int, from time.Time) ([]models.Birthday, error) {
	var rows *sql.Rows
	var err error
	fromMonth := int(from.Month())
	if fromMonth+months > 12 {
		rows, err = br.db.Query(`SELECT * FROM birthdays WHERE (EXTRACT(MONTH FROM date) > $1) OR (EXTRACT(MONTH FROM date) <= $2)`, fromMonth, fromMonth+months-12)
	} else {
		rows, err = br.db.Query(`SELECT * FROM birthday WHERE (EXTRACT(MONTH FROM date) > $1) AND (EXTRACT(MONTH FROM date) <= $2)`, fromMonth, fromMonth+months)
	}
	return scanBirthdayRows(rows, err)
}

func (br *birthdayDAOSQL) GetBirthdaysOfMonth(month int) ([]models.Birthday, error) {
	rows, err := br.db.Query(`SELECT * FROM birthdays WHERE EXTRACT(MONTH FROM date) = $1`, month)
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

func (br *birthdayDAOSQL) GetBirthdayByUserId(userId string) (models.Birthday, error) {
	var result models.Birthday
	err := br.db.QueryRow(`SELECT * FROM birthdays WHERE user_id = $1`, userId).Scan(&result.Id, &result.UserId, &result.Date)
	return result, err
}

func (br *birthdayDAOSQL) AddBirthday(userId string, date time.Time) (*models.Birthday, error) {
	result := &models.Birthday{
		UserId: userId,
		Date:   date,
	}
	err := br.db.QueryRow(`INSERT INTO birthdays (user_id, date) VALUES ($1, $2) RETURNING Id`, userId, date).Scan(&result.Id)
	return result, err
}

func (br *birthdayDAOSQL) RemoveBirthday(userId string) error {
	return br.db.QueryRow(`DELETE FROM birthdays WHERE user_id = $1`, userId).Err()
}
