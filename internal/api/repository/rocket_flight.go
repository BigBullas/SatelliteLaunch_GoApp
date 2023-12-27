package repository

import (
	"RIP_lab1/internal/models"
	"time"
)

func (r *Repository) GetRocketFlightList(formDateStart time.Time, formDateEnd time.Time, status string) ([]models.RocketFlight, error) {
	var rocketFlights []models.RocketFlight

	if status != "" {
		if formDateStart.IsZero() {
			if formDateEnd.IsZero() {
				// фильтрация только по статусу
				res := r.db.Where("status = ?", status).Find(&rocketFlights)
				return rocketFlights, res.Error
			}

			// фильтрация по статусу и formDateEnd
			res := r.db.Where("status = ?", status).Where("formed_at < ?", formDateEnd).Find(&rocketFlights)
			return rocketFlights, res.Error
		}

		// фильтрация по статусу и formDateStart
		if formDateEnd.IsZero() {
			res := r.db.Where("status = ?", status).Where("formed_at > ?", formDateStart).
				Find(&rocketFlights)
			return rocketFlights, res.Error
		}

		// фильтрация по статусу, formDateStart и formDateEnd
		res := r.db.Where("status = ?", status).Where("formed_at BETWEEN ? AND ?", formDateStart, formDateEnd).Find(&rocketFlights)
		return rocketFlights, res.Error
	}

	if formDateStart.IsZero() {
		if formDateEnd.IsZero() {
			// без фильтрации
			res := r.db.Where("status IN (?)", []string{"formed", "completed", "rejected"}).Find(&rocketFlights)
			return rocketFlights, res.Error
		}

		// фильтрация по formDateEnd
		res := r.db.Where("status IN (?)", []string{"formed", "completed", "rejected"}).Where("formed_at < ?", formDateEnd).Find(&rocketFlights)
		return rocketFlights, res.Error
	}

	if formDateEnd.IsZero() {
		// фильтрация по formDateStart
		res := r.db.Where("status IN (?)", []string{"formed", "completed", "rejected"}).Where("formed_at > ?", formDateStart).Find(&rocketFlights)
		return rocketFlights, res.Error
	}

	//фильтрация по formDateStart и formDateEnd
	res := r.db.Where("status IN (?)", []string{"formed", "completed", "rejected"}).
		Where("formed_at BETWEEN ? AND ?", formDateStart, formDateEnd).Find(&rocketFlights)
	return rocketFlights, res.Error
}

func (r *Repository) GetRocketFlightById(flightId int) (models.RocketFlight, []models.FlightRequest, error) {
	var rocketFlight models.RocketFlight
	var flightRequests []models.FlightRequest

	//информация по данному полёту
	result := r.db.First(&rocketFlight, "flight_id =?", flightId)
	if result.Error != nil {
		// log.Println("Ошибка при получении данного полёта")
		return models.RocketFlight{}, []models.FlightRequest{}, result.Error
	}

	//заявки на полёт КА, принятые на данный полёт
	result = r.db.Table("flights_flight_requests").Select("flight_requests.*").
		Joins("JOIN flight_requests ON flights_flight_requests.request_id = flight_requests.request_id").
		Where("flights_flight_requests.flight_id = ?", flightId).Find(&flightRequests)
	if result.Error != nil {
		// log.Println("Ошибка при получении заявок на полёт КА по данному полёту")
		return models.RocketFlight{}, []models.FlightRequest{}, result.Error
	}

	return rocketFlight, flightRequests, nil
}
