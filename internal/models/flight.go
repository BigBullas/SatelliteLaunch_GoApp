package models

import "time"

type RocketFlight struct {
	FlightId       int       `gorm:"primarykey" json:"flightId,omitempty"`
	CreatorId      int       `json:"creatorId,omitempty"`
	CreatorLogin   string    `json:"creatorLogin,omitempty" gorm:"-"`
	ModeratorId    int       `json:"moderatorId,omitempty"`
	ModeratorLogin string    `json:"moderatorLogin,omitempty" gorm:"-"`
	Status         string    `json:"status,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
	FormedAt       time.Time `json:"formedAtomitempty"`
	ConfirmedAt    time.Time `json:"confirmedAt,omitempty"`
	FlightDate     time.Time `json:"flightDate,omitempty"`
	LoadCapacity   int       `json:"loadCapacity,omitempty"`
	Price          float64   `json:"price,omitempty"`
	Title          string    `json:"title,omitempty"`
	PlaceNumber    int       `json:"placeNumber,omitempty"`
}
