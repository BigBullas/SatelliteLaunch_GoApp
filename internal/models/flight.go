package models

import "time"

type RocketFlight struct {
	FlightId    int `gorm:"primarykey"`
	CreatorId   int
	ModeratorId int
	Status      string
	CreatedAt   time.Time
	FormedAt    time.Time
	ConfirmedAt time.Time
	FlightDate  time.Time
	Payload     int
	Price       float64
	Title       string
	PlaceNumber int
}

type ShortRocketFlight struct {
	CreatorId int
	RequestId int
}
