package api

import "RIP_lab1/internal/models"

type Repo interface {
	GetRequestForFlightList(substring string) ([]models.FlightRequest, error)
	GetCardRequestForFlightById(cardId int) (models.FlightRequest, error)
	CreateNewRequestForFlight(models.FlightRequest) error
	ChangeRequestForFlight(models.FlightRequest) error
	DeleteRequestForFlightById(cardId int) error
}
