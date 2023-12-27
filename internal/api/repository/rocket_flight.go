package repository

import (
	"RIP_lab1/internal/models"
	"fmt"
	"time"

	"gorm.io/gorm"
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

func (r *Repository) GetRocketFlightDraft(userId int) (int, error) {
	var rocketFlight models.RocketFlight
	err := r.db.First(&rocketFlight, "creator_id = ? and status = 'draft'", userId)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		return 0, err.Error
	}

	return rocketFlight.FlightId, nil
}

func (r *Repository) GetRocketFlightById(flightId int) (models.RocketFlightDetailed, []models.FlightRequest, error) {
	var rocketFlight models.RocketFlight
	// var rocketFlightDetailed models.RocketFlightDetailed
	var flightRequests []models.FlightRequest

	//информация по данному полёту
	result := r.db.First(&rocketFlight, "flight_id =?", flightId)
	if result.Error != nil {
		// log.Println("Ошибка при получении данного полёта")
		return models.RocketFlightDetailed{}, []models.FlightRequest{}, result.Error
	}

	var creator models.UserShort
	result = r.db.Table("users").Select("login").Where("user_id = ?", 1).First(&creator)
	if result.Error != nil {
		// log.Println("Ошибка при получении данного полёта")
		return models.RocketFlightDetailed{}, []models.FlightRequest{}, result.Error
	}

	var moderator models.UserShort
	result = r.db.Table("users").Select("login").Where("user_id = ?", 2).First(&moderator)
	if result.Error != nil {
		// log.Println("Ошибка при получении данного полёта")
		return models.RocketFlightDetailed{}, []models.FlightRequest{}, result.Error
	}

	rocketFlightDetailed := models.RocketFlightDetailed{
		FlightId:       rocketFlight.FlightId,
		CreatorLogin:   creator.Login,
		ModeratorLogin: moderator.Login,
		Status:         rocketFlight.Status,
		CreatedAt:      rocketFlight.CreatedAt,
		FormedAt:       rocketFlight.FormedAt,
		ConfirmedAt:    rocketFlight.ConfirmedAt,
		FlightDate:     rocketFlight.FlightDate,
		Payload:        rocketFlight.Payload,
		Price:          rocketFlight.Price,
		Title:          rocketFlight.Title,
		PlaceNumber:    rocketFlight.PlaceNumber,
	}

	//заявки на полёт КА, принятые на данный полёт
	result = r.db.Table("flights_flight_requests").Select("flight_requests.*").
		Joins("JOIN flight_requests ON flights_flight_requests.request_id = flight_requests.request_id").
		Where("flights_flight_requests.flight_id = ?", flightId).Find(&flightRequests)
	if result.Error != nil {
		// log.Println("Ошибка при получении заявок на полёт КА по данному полёту")
		return models.RocketFlightDetailed{}, []models.FlightRequest{}, result.Error
	}

	return rocketFlightDetailed, flightRequests, nil
}

func (r *Repository) ChangeRocketFlight(changedRocketFlight models.RocketFlightChangeable) error {
	var oldRocketFlight models.RocketFlight
	result := r.db.Table("rocket_flights").Where("creator_id = ? AND status = ?", changedRocketFlight.CreatorId, "draft").Find(&oldRocketFlight)
	if result.Error != nil {
		return result.Error
	}

	if !changedRocketFlight.FlightDate.IsZero() {
		oldRocketFlight.FlightDate = changedRocketFlight.FlightDate
	}

	if changedRocketFlight.Payload != 0 {
		oldRocketFlight.Payload = changedRocketFlight.Payload
	}

	if changedRocketFlight.Price != 0.0 {
		oldRocketFlight.Price = changedRocketFlight.Price
	}

	if changedRocketFlight.Title != "" {
		oldRocketFlight.Title = changedRocketFlight.Title
	}

	if changedRocketFlight.PlaceNumber != 0 {
		oldRocketFlight.PlaceNumber = changedRocketFlight.PlaceNumber
	}

	result = r.db.Save(oldRocketFlight)
	return result.Error
}

func (r *Repository) FormRocketFlight(newRocketStatus models.RocketFlight) error {
	var rocketFlight models.RocketFlight
	err := r.db.First(&rocketFlight, "creator_id = ? and status = 'draft'", newRocketStatus.CreatorId)
	if err.Error != nil {
		return err.Error
	}

	rocketFlight.Status = newRocketStatus.Status
	rocketFlight.FormedAt = time.Now()

	result := r.db.Save(&rocketFlight)

	return result.Error
}

func (r *Repository) ResponceRocketFlight(newFlightStatus models.RocketFlight) error {
	var rocketFlight models.RocketFlight

	err := r.db.First(&rocketFlight, "flight_id = ? and status = 'formed'", newFlightStatus.FlightId)
	if err.Error != nil && err.Error.Error() == "record not found" {
		return fmt.Errorf("Такой заявки-черновика на полёт ракеты-носителя нет")
	}
	if err.Error != nil {
		return err.Error
	}

	rocketFlight.Status = newFlightStatus.Status
	rocketFlight.ModeratorId = newFlightStatus.ModeratorId

	res := r.db.Save(&rocketFlight)
	return res.Error
}

func (r *Repository) DeleteRocketFlight(userId int) error {
	var rocketFlight models.RocketFlight
	result := r.db.First(&rocketFlight, "creator_id =? and status = 'draft'", userId)
	if result.Error != nil {
		return result.Error
	}

	rocketFlight.Status = "deleted"
	result = r.db.Save(rocketFlight)
	return result.Error
}
