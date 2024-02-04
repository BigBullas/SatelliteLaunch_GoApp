package models

import "time"

type RocketFlight struct {
	FlightId       int       `gorm:"primarykey" json:"flight_id,omitempty"`
	CreatorId      int       `json:"creator_id,omitempty"`
	CreatorLogin   string    `json:"creator_login,omitempty" gorm:"-"`
	ModeratorId    int       `json:"moderator_id,omitempty"`
	ModeratorLogin string    `json:"moderator_login,omitempty" gorm:"-"`
	Status         string    `json:"status,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	FormedAt       time.Time `json:"formed_at,omitempty"`
	ConfirmedAt    time.Time `json:"confirmed_at,omitempty"`
	FlightDate     time.Time `json:"flight_date,omitempty"`
	LoadCapacity   float64   `json:"load_capacity,omitempty"`
	Price          float64   `json:"price,omitempty"`
	Title          string    `json:"title,omitempty"`
	PlaceNumber    int       `json:"place_number,omitempty"`
}

type FlightAsync struct {
	Id              int    `json:"flight_id"`
	CalculatedPrice int    `json:"calculated_price"`
	Token           string `json:"token"`
}
