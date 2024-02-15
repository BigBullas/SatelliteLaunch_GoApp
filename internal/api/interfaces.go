package api

import "RIP_lab1/internal/models"

type Repo interface {
	GetRequestForDeliveryList(substring string) (int64, []models.FlightRequest, error)
	AddFlightRequestToFlight(creatorId int, requestId int) (error)
	GetCardRequestForDeliveryByID(cardId int) (models.FlightRequest, error)
	DeleteRequestForDeliveryById(cardId int) error
}
