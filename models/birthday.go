package models

import "time"

type Birthday struct {
	Id     uint32
	UserId string
	Date   time.Time
}
