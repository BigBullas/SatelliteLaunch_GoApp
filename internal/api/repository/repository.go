package repository

import (
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"RIP_lab1/internal/models"
)

type Repository struct {
	db *gorm.DB
}

const draftId = 1

func NewRepo(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.Payload{})
	if err != nil {
		panic("Миграция БД не удалась")
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetRequestForDeliveryList(substring string) (int64, []models.Payload, error) {
	var request_for_delivery []models.Payload
	var count int64

	err := r.db.Table("flights_payloads").Where("flight_id = ?", draftId).Count(&count).Error
	if err != nil {
		return 0, request_for_delivery, err
	}

	r.db.Where("title ILIKE ?", "%"+substring+"%").Find(&request_for_delivery, "is_available = ?", true)
	return count, request_for_delivery, nil
}

func (r *Repository) AddFlightRequestToFlight(requestId int) error {
	payloadFlight := models.FlightsPayloads{
		FlightId:        draftId,
		PayloadId:       requestId,
		CountSatellites: 1,
	}
	res := r.db.Create(&payloadFlight)

	if res.Error != nil && res.Error.Error() == "ERROR: duplicate key value violates unique constraint \"flights_payloads_pkey\" (SQLSTATE 23505)" {
		return errors.New("Данная полезная нагрузка уже добавлена в планируемый полёт")
	}

	return res.Error
}

func (r *Repository) GetCardRequestForDeliveryByID(payloadId int) (models.Payload, error) {
	var card models.Payload

	r.db.Where("payload_id = ?", payloadId).Find(&card, "is_available = ?", true)
	return card, nil

}
func (r *Repository) DeleteRequestForDeliveryById(cardId int) error {
	err := r.db.Exec("UPDATE payloads SET is_available=false WHERE payload_id = ?", cardId).Error

	if err != nil {
		return err
	}
	return nil
}
