package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"RIP_lab1/internal/models"
)

type Repository struct {
	db *gorm.DB
}

func NewRepo(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.FlightRequest{})
	if err != nil {
		panic("Миграция БД не удалась")
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetRequestForDeliveryList(substring string) ([]models.FlightRequest, error) {
	var request_for_delivery []models.FlightRequest

	r.db.Where("title ILIKE ?", "%"+substring+"%").Find(&request_for_delivery, "is_available = ?", true)
	return request_for_delivery, nil
}

func (r *Repository) GetCardRequestForDeliveryByID(cardId int) (models.FlightRequest, error) {
	var card models.FlightRequest

	r.db.Where("request_id = ?", cardId).Find(&card, "is_available = ?", true)
	return card, nil
}

func (r *Repository) DeleteRequestForDeliveryById(cardId int) error {
	r.db.Where("request_id = ? AND is_available = ?", cardId, true).Update("is_available", false)
	return nil
}
