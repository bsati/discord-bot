package services

import (
	"database/sql"

	"github.com/bsati/discord-bot/daos"
)

type ServiceRegistry struct {
	BirthdayService
}

func InitServices(db *sql.DB) *ServiceRegistry {
	birthdayDAO := daos.NewBirthdayDAO(db)

	return &ServiceRegistry{
		BirthdayService: NewBirthdayService(birthdayDAO),
	}
}
