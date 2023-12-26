package repository

import (
	"time"

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

func (r *Repository) GetRequestForFlightList(substring string) ([]models.FlightRequest, error) {
	var request_for_delivery []models.FlightRequest

	r.db.Where("title ILIKE ?", "%"+substring+"%").Find(&request_for_delivery, "is_available = ?", true)
	return request_for_delivery, nil
}

func (r *Repository) GetCardRequestForFlightById(cardId int) (models.FlightRequest, error) {
	var card models.FlightRequest

	r.db.Where("request_id = ?", cardId).Find(&card, "is_available = ?", true)
	return card, nil
}

func (r *Repository) CreateNewRequestForFlight(newFlightRequest models.FlightRequest) error {
	result := r.db.Create(&newFlightRequest)
	return result.Error
}

func (r *Repository) ChangeRequestForFlight(changedFlightRequest models.FlightRequest) error {
	var oldFlightRequest models.FlightRequest
	result := r.db.First(&oldFlightRequest, "request_id =?", changedFlightRequest.RequestId)
	if result.Error != nil {
		return result.Error
	}

	if changedFlightRequest.Title != "" {
		oldFlightRequest.Title = changedFlightRequest.Title
	}

	if changedFlightRequest.Description != "" {
		oldFlightRequest.Description = changedFlightRequest.Description
	}

	if changedFlightRequest.DetailedDesc != "" {
		oldFlightRequest.DetailedDesc = changedFlightRequest.DetailedDesc
	}

	if changedFlightRequest.DesiredPrice != 0.0 {
		oldFlightRequest.DesiredPrice = changedFlightRequest.DesiredPrice
	}

	if changedFlightRequest.LoadCapacity != 0.0 {
		oldFlightRequest.LoadCapacity = changedFlightRequest.LoadCapacity
	}

	if changedFlightRequest.ImgURL != "" {
		oldFlightRequest.ImgURL = changedFlightRequest.ImgURL
	}

	if changedFlightRequest.FlightDateStart != (time.Time{}) {
		oldFlightRequest.FlightDateStart = changedFlightRequest.FlightDateStart
	}

	if changedFlightRequest.FlightDateEnd != (time.Time{}) {
		oldFlightRequest.FlightDateEnd = changedFlightRequest.FlightDateEnd
	}

	result = r.db.Save(oldFlightRequest)
	return result.Error
}

func (r *Repository) DeleteRequestForFlightById(cardId int) error {
	err := r.db.Exec("UPDATE flight_requests SET is_available=false WHERE request_id = ?", cardId).Error
	if err != nil {
		return err
	}
	return nil
}
