package daos

import (
	"database/sql"
	"time"

	"github.com/bsati/discord-bot/models"
)

type BirthdayDAO interface {
	GetBirthdays() ([]models.Birthday, error)
	GetBirthdaysByDate(date time.Time) ([]models.Birthday, error)
	GetBirthdayByUserId(userId string) (models.Birthday, error)
	AddBirthday(userId string, date time.Time) (*models.Birthday, error)
	RemoveBirthday(userId string) error
}

type birthdayDAOSQL struct {
	db *sql.DB
}

func NewBirthdayDAO(db *sql.DB) BirthdayDAO {
	return &birthdayDAOSQL{db: db}
}

func (br *birthdayDAOSQL) GetBirthdays() ([]models.Birthday, error) {
	rows, err := br.db.Query(`SELECT * FROM Birthdays`)
	return scanBirthdayRows(rows, err)
}

func (br *birthdayDAOSQL) GetBirthdaysByDate(date time.Time) ([]models.Birthday, error) {
	rows, err := br.db.Query(`SELECT * FROM Birthdays WHERE Date = ?`, date)
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
	err := br.db.QueryRow(`SELECT * FROM Birthdays WHERE UserId = ?`, userId).Scan(&result.Id, &result.UserId, &result.Date)
	return result, err
}

func (br *birthdayDAOSQL) AddBirthday(userId string, date time.Time) (*models.Birthday, error) {
	result := &models.Birthday{
		UserId: userId,
		Date:   date,
	}
	err := br.db.QueryRow(`INSERT INTO Birthdays(UserId, Date) VALUES (?, ?) RETURNING Id`, userId, date).Scan(result.Id)
	return result, err
}

func (br *birthdayDAOSQL) RemoveBirthday(userId string) error {
	return br.db.QueryRow(`DELETE FROM Birthdays WHERE UserId = ?`, userId).Err()
}
