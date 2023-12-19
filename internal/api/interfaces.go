package api

import "RIP_lab1/internal/models"

type Repo interface {
	GetStarsByNameFilter(substring string) ([]models.Product, error)
	GetStarByID(threatId int) (models.Product, error)
	DeleteStarById(starId int) error

	GetProductByID(id uint) (*models.Product, error)
}