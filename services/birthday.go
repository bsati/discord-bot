package services

import (
	"time"

	"github.com/bsati/discord-bot/daos"
	"github.com/bsati/discord-bot/models"
)

type BirthdayService interface {
	GetBirthdays() ([]models.Birthday, error)
	AddBirthday(userId string, date time.Time) (*models.Birthday, error)
	RemoveBirthday(userId string) error
}

type birthdayService struct {
	birthdayDAO daos.BirthdayDAO
}

func NewBirthdayService(birthdayDAO daos.BirthdayDAO) BirthdayService {
	return &birthdayService{birthdayDAO: birthdayDAO}
}

func (bs *birthdayService) GetBirthdays() ([]models.Birthday, error) {
	return bs.birthdayDAO.GetBirthdays()
}

func (bs *birthdayService) AddBirthday(userId string, date time.Time) (*models.Birthday, error) {
	return bs.birthdayDAO.AddBirthday(userId, date)
}

func (bs *birthdayService) RemoveBirthday(userId string) error {
	return bs.birthdayDAO.RemoveBirthday(userId)
}
