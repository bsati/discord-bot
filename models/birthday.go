package models

import "time"

// Birthday represents a user's birthday
type Birthday struct {
	Id     uint32
	UserId string
	Date   time.Time
}
