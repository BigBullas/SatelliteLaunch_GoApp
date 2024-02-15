package api

import "RIP_lab1/internal/models"

type Repo interface {
	GetRequestForDeliveryList(substring string) (int64, []models.Payload, error)
	AddFlightRequestToFlight(requestId int) error
	GetCardRequestForDeliveryByID(cardId int) (models.Payload, error)
	DeleteRequestForDeliveryById(cardId int) error
}
