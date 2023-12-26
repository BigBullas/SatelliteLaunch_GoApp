package repository

import (
	"errors"
	"time"

	"RIP_lab1/internal/models"
)

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

func (r *Repository) AddFlightRequestToFlight(shortFlight models.ShortRocketFlight) error {
	var rocketFlight models.RocketFlight

	r.db.Where("creator_id = ?", shortFlight.CreatorId).Where("status = ?", "draft").First(&rocketFlight)

	// log.Println(rocketFlight)

	if rocketFlight.FlightId == 0 {
		newRocketFlights := models.RocketFlight{
			CreatorId:   shortFlight.CreatorId,
			ModeratorId: 2,
			Status:      "draft",
			CreatedAt:   time.Now(),
		}
		res := r.db.Create(&newRocketFlights)
		if res.Error != nil {
			return res.Error
		}
		rocketFlight = newRocketFlights
	}

	flightRequestFlight := models.FlightsFlightRequest{
		FlightId:  rocketFlight.FlightId,
		RequestId: shortFlight.RequestId,
	}

	res := r.db.Create(&flightRequestFlight)
	if res.Error != nil && res.Error.Error() == "ERROR: duplicate key value violates unique constraint \"flights_flight_requests_pkey\" (SQLSTATE 23505)" {
		return errors.New("Данная заявка на доставку КА уже добавлена в планируемый полёт")
	}

	return res.Error
}
