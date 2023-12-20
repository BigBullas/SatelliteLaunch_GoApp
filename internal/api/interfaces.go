package api

import "RIP_lab1/internal/models"

type Repo interface {
	GetRequestForDeliveryList(substring string) ([]models.FlightRequest, error)
	GetCardRequestForDeliveryByID(cardId int) (models.FlightRequest, error)
	DeleteRequestForDeliveryById(cardId int) error
}
