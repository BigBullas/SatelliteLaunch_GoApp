package api

import "RIP_lab1/internal/models"

type Repo interface {
	GetRequestForDeliveryList(substring string) ([]models.RequestForDelivery, error)
	GetCardRequestForDeliveryByID(cardId int) (models.RequestForDelivery, error)
	DeleteRequestForDeliveryById(cardId int) error
}
