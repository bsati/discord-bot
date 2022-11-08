package daos

import "database/sql"

type DAO struct {
	BirthdayDAO
}

func NewDAO(db *sql.DB) *DAO {
	return &DAO{
		BirthdayDAO: &birthdayDAOSQL{db},
	}
}
