package api

import (
	"RIP_lab1/internal/models"
	"time"
)

type Repo interface {
	GetRequestForFlightList(substring string) ([]models.FlightRequest, error)
	GetCardRequestForFlightById(cardId int) (models.FlightRequest, error)
	CreateNewRequestForFlight(models.FlightRequest) error
	ChangeRequestForFlight(models.FlightRequest) error
	DeleteRequestForFlightById(cardId int) error

	AddFlightRequestToFlight(models.RocketFlightShort) error

	GetRocketFlightList(time.Time, time.Time, string) ([]models.RocketFlight, error)
	GetRocketFlightById(int) (models.RocketFlightDetailed, []models.FlightRequest, error)
	ChangeRocketFlight(models.RocketFlightChangeable) error
}
