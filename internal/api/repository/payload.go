package repository

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"RIP_lab1/internal/models"
)

func (r *Repository) GetPayloadList(substring string, loadCapacityStart string,
	loadCapacityEnd string, flightDateStart string, flightDateEnd string) ([]models.Payload, error) {
	var request_for_delivery []models.Payload
	queryParametrs := "is_available = true"

	if loadCapacityStart != "" {
		queryParametrs = fmt.Sprintf(queryParametrs+" AND load_capacity >= '%s'", loadCapacityStart)
	}

	if loadCapacityEnd != "" {
		queryParametrs = fmt.Sprintf(queryParametrs+" AND load_capacity <= '%s'", loadCapacityEnd)
	}

	if flightDateStart != "" {
		queryParametrs = fmt.Sprintf(queryParametrs+" AND flight_date_end >= '%s'", flightDateStart)
	}

	if flightDateEnd != "" {
		queryParametrs = fmt.Sprintf(queryParametrs+" AND flight_date_start <= '%s'", flightDateEnd)
	}

	r.db.Where("title ILIKE ?", "%"+substring+"%").Find(&request_for_delivery, queryParametrs)
	return request_for_delivery, nil
}

func (r *Repository) GetCardPayloadById(cardId int) (models.Payload, error) {
	var card models.Payload

	r.db.Where("payload_id = ?", cardId).Find(&card, "is_available = ?", true)
	return card, nil
}

func (r *Repository) CreateNewPayload(newPayload models.Payload) error {
	result := r.db.Create(&newPayload)
	return result.Error
}

func (r *Repository) ChangePayload(changedPayload models.Payload) error {
	var oldPayload models.Payload
	result := r.db.First(&oldPayload, "payload_id =?", changedPayload.PayloadId)
	if result.Error != nil {
		return result.Error
	}

	if changedPayload.Title != "" {
		oldPayload.Title = changedPayload.Title
	}

	if changedPayload.Description != "" {
		oldPayload.Description = changedPayload.Description
	}

	if changedPayload.DetailedDesc != "" {
		oldPayload.DetailedDesc = changedPayload.DetailedDesc
	}

	if changedPayload.DesiredPrice != 0.0 {
		oldPayload.DesiredPrice = changedPayload.DesiredPrice
	}

	if changedPayload.LoadCapacity != 0.0 {
		oldPayload.LoadCapacity = changedPayload.LoadCapacity
	}

	if changedPayload.ImgURL != "" {
		oldPayload.ImgURL = changedPayload.ImgURL
	}

	if changedPayload.FlightDateStart != (time.Time{}) {
		oldPayload.FlightDateStart = changedPayload.FlightDateStart
	}

	if changedPayload.FlightDateEnd != (time.Time{}) {
		oldPayload.FlightDateEnd = changedPayload.FlightDateEnd
	}

	result = r.db.Save(oldPayload)
	return result.Error
}

func (r *Repository) DeletePayloadById(cardId int) error {
	err := r.db.Exec("UPDATE payloads SET is_available=false WHERE payload_id = ?", cardId).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddPayloadToFlight(creatorId int, requestId int) error {
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
		log.Println("flightId", rocketFlight.FlightId)
		log.Println("payloadId", requestId)
	}

	payloadFlight := models.FlightsPayload{
		FlightId:        rocketFlight.FlightId,
		PayloadId:       requestId,
		CountSatellites: 1,
	}

	res := r.db.Create(&payloadFlight)
	if res.Error != nil && res.Error.Error() == "ERROR: duplicate key value violates unique constraint \"flights_payloads_pkey\" (SQLSTATE 23505)" {
		return errors.New("Данная полезная нагрузка уже добавлена в планируемый полёт")
	}

	return res.Error
}

func (r *Repository) DeletePayloadFromFlight(userId int, requestId int) error {
	var rocketFlight models.RocketFlight
	r.db.Where("creator_id = ? and status = 'draft'", userId).First(&rocketFlight)

	if rocketFlight.FlightId == 0 {
		return errors.New("Нет заявки-черновика на полёт ракеты-носителя")
	}

	var flightsPayload models.FlightsPayload
	err := r.db.Where("flight_id = ? AND payload_id = ?", rocketFlight.FlightId, requestId).First(&flightsPayload).Error
	if err != nil {
		return errors.New("Такой полезной нагрузки нет в данном планируемом полёте")
	}

	err = r.db.Where("flight_id = ? AND payload_id = ?", rocketFlight.FlightId, requestId).Delete(models.FlightsPayload{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ChangeCountFlightsPayload(userId int, requestId int, count int) error {
	var rocketFlight models.RocketFlight
	r.db.Where("creator_id = ? and status = 'draft'", userId).First(&rocketFlight)

	if rocketFlight.FlightId == 0 {
		return errors.New("Нет заявки-черновика на полёт ракеты-носителя")
	}

	var flightsPayload models.FlightsPayload
	err := r.db.Where("flight_id = ? AND payload_id = ?", rocketFlight.FlightId, requestId).First(&flightsPayload).Error
	if err != nil {
		return errors.New("Такой полезной нагрузки нет в данном планируемом полёте")
	}

	flightsPayload.CountSatellites += count
	if flightsPayload.CountSatellites < 1 {
		return errors.New("Количество данных полезных нагрузок будет меньше одного после выполнения данной операции")
	}

	err = r.db.Where("flight_id = ? AND payload_id = ?", rocketFlight.FlightId, requestId).Save(flightsPayload).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetPayloadImageUrl(requestId int) string {
	payload := models.Payload{}

	r.db.First(&payload, "payload_id = ?", strconv.Itoa(requestId))
	return payload.ImgURL
}
