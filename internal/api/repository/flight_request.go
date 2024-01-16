package repository

import (
	"errors"
	"strconv"
	"time"

	"RIP_lab1/internal/models"
)

func (r *Repository) GetRequestForFlightList(substring string) ([]models.Payload, error) {
	var request_for_delivery []models.Payload

	r.db.Where("title ILIKE ?", "%"+substring+"%").Find(&request_for_delivery, "is_available = ?", true)
	return request_for_delivery, nil
}

func (r *Repository) GetCardRequestForFlightById(cardId int) (models.Payload, error) {
	var card models.Payload

	r.db.Where("request_id = ?", cardId).Find(&card, "is_available = ?", true)
	return card, nil
}

func (r *Repository) CreateNewRequestForFlight(newFlightRequest models.Payload) error {
	result := r.db.Create(&newFlightRequest)
	return result.Error
}

func (r *Repository) ChangeRequestForFlight(changedFlightRequest models.Payload) error {
	var oldFlightRequest models.Payload
	result := r.db.First(&oldFlightRequest, "request_id =?", changedFlightRequest.PayloadId)
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

func (r *Repository) AddFlightRequestToFlight(creatorId int, requestId int) error {
	var rocketFlight models.RocketFlight

	r.db.Where("creator_id = ?", creatorId).Where("status = ?", "draft").First(&rocketFlight)

	// log.Println(rocketFlight)

	if rocketFlight.FlightId == 0 {
		newRocketFlights := models.RocketFlight{
			CreatorId:   creatorId,
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

	flightRequestFlight := models.FlightsPayload{
		FlightId:        rocketFlight.FlightId,
		PayloadId:       requestId,
		CountSatellites: 1,
	}

	res := r.db.Create(&flightRequestFlight)
	if res.Error != nil && res.Error.Error() == "ERROR: duplicate key value violates unique constraint \"flights_flight_requests_pkey\" (SQLSTATE 23505)" {
		return errors.New("Данная заявка на доставку КА уже добавлена в планируемый полёт")
	}

	return res.Error
}

func (r *Repository) DeleteRequestFromFlight(userId int, requestId int) error {
	var rocketFlight models.RocketFlight
	r.db.Where("creator_id = ? and status = 'draft'", userId).First(&rocketFlight)

	if rocketFlight.FlightId == 0 {
		return errors.New("Нет заявки-черновика на полёт ракеты-носителя")
	}

	var flightsFlightRequest models.FlightsPayload
	err := r.db.Where("flight_id = ? AND request_id = ?", rocketFlight.FlightId, requestId).First(&flightsFlightRequest).Error
	if err != nil {
		return errors.New("Такой заявки нет в данном планируемом полёте")
	}

	err = r.db.Where("flight_id = ? AND request_id = ?", rocketFlight.FlightId, requestId).Delete(models.FlightsPayload{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ChangeCountFlightsFlightRequest(userId int, requestId int, count int) error {
	var rocketFlight models.RocketFlight
	r.db.Where("creator_id = ? and status = 'draft'", userId).First(&rocketFlight)

	if rocketFlight.FlightId == 0 {
		return errors.New("Нет заявки-черновика на полёт ракеты-носителя")
	}

	var flightsFlightRequest models.FlightsPayload
	err := r.db.Where("flight_id = ? AND request_id = ?", rocketFlight.FlightId, requestId).First(&flightsFlightRequest).Error
	if err != nil {
		return errors.New("Такой заявки нет в данном планируемом полёте")
	}

	flightsFlightRequest.CountSatellites += count
	if flightsFlightRequest.CountSatellites < 1 {
		return errors.New("Количество заявок на полёт данного КА меньше, чем изменение количества спутников")
	}

	err = r.db.Where("flight_id = ? AND request_id = ?", rocketFlight.FlightId, requestId).Save(flightsFlightRequest).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetFlightRequestImageUrl(requestId int) string {
	flightRequest := models.Payload{}

	r.db.First(&flightRequest, "request_id = ?", strconv.Itoa(requestId))
	return flightRequest.ImgURL
}
