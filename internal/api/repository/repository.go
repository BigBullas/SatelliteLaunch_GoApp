package repository

import (
	// "log"

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
	err = db.AutoMigrate(&models.FlightRequest{})
	if err != nil {
		panic("Миграция БД не удалась")
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetRequestForDeliveryList(substring string) (int64, []models.FlightRequest, error) {
	var request_for_delivery []models.FlightRequest

	var count int64

	err := r.db.Table("flights_flight_requests").Where("flight_id = ?", 0).Count(&count).Error

	if err != nil {
		return 0, request_for_delivery, err
	}

	r.db.Where("title ILIKE ?", "%"+substring+"%").Find(&request_for_delivery, "is_available = ?", true)

	return 0, request_for_delivery, nil

}
func (r *Repository) AddFlightRequestToFlight(creatorId int, requestId int) error {
	// var rocketFlight models.RocketFlight

	// r.db.Where("creator_id = ?", creatorId).Where("status = ?", "draft").First(&rocketFlight)

	// log.Println(rocketFlight)

	// if rocketFlight.FlightId == 0 {

	// 	newRocketFlights := models.RocketFlight{
	// 		CreatorId:   creatorId,
	// 		ModeratorId: 2,
	// 		Status:      "draft",
	// 		CreatedAt:   time.Now(),
	// 	}
	// 	res := r.db.Create(&newRocketFlights)
	// 	if res.Error != nil {
	// 		return res.Error
	// 	}
	// 	rocketFlight = newRocketFlights
	// }
	payloadFlight := models.FlightsFlightRequests{
		FlightId: draftId,

		PayloadId:       requestId,
		CountSatellites: 1,
	}
	res := r.db.Create(&payloadFlight)

	if res.Error != nil && res.Error.Error() == "ERROR: duplicate key value violates unique constraint \"flights_payloads_pkey\" (SQLSTATE 23505)" {
		return errors.New("Данная полезная нагрузка уже добавлена в планируемый полёт")
	}

	return res.Error
}

func (r *Repository) GetCardRequestForDeliveryByID(cardId int) (models.FlightRequest, error) {
	var card models.FlightRequest

	r.db.Where("request_id = ?", cardId).Find(&card, "is_available = ?", true)
	return card, nil

}
func (r *Repository) DeleteRequestForDeliveryById(cardId int) error {
	err := r.db.Exec("UPDATE flight_requests SET is_available=false WHERE request_id = ?", cardId).Error

	if err != nil {
		return err
	}
	return nil
}
