package repository

import (
	"RIP_lab1/internal/models"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func (r *Repository) GetUserByUserId(user_id int) (models.User, error) {
	user := models.User{}
	result := r.db.Table("users").Select("login").Where("user_id = ?", user_id).First(&user)
	return user, result.Error
}

func (r *Repository) GetLoginsForFlights(rocketFlights []models.RocketFlight) ([]models.RocketFlight, error) {
	var creator models.User
	var moderator models.User
	var err error

	for i := range rocketFlights {
		creator, err = r.GetUserByUserId(rocketFlights[i].CreatorId)
		moderator, err = r.GetUserByUserId(rocketFlights[i].ModeratorId)
		if err != nil {
			return rocketFlights, err
		}
		rocketFlights[i].CreatorLogin = creator.Login
		rocketFlights[i].ModeratorLogin = moderator.Login
	}

	return rocketFlights, nil
}

func (r *Repository) GetRocketFlightList(formDateStart time.Time, formDateEnd time.Time, status string) ([]models.RocketFlight, error) {
	var rocketFlights []models.RocketFlight
	var err error

	if status != "" {
		if formDateStart.IsZero() {
			if formDateEnd.IsZero() {
				// фильтрация только по статусу
				res := r.db.Where("status = ?", status).Find(&rocketFlights)
				if res.Error != nil {
					return rocketFlights, res.Error
				}

				rocketFlights, err = r.GetLoginsForFlights(rocketFlights)
				if err != nil {
					return rocketFlights, err
				}

				return rocketFlights, nil
			}

			// фильтрация по статусу и formDateEnd
			res := r.db.Where("status = ?", status).Where("formed_at < ?", formDateEnd).Find(&rocketFlights)
			if res.Error != nil {
				return rocketFlights, res.Error
			}

			rocketFlights, err = r.GetLoginsForFlights(rocketFlights)
			if err != nil {
				return rocketFlights, err
			}

			return rocketFlights, nil
		}

		// фильтрация по статусу и formDateStart
		if formDateEnd.IsZero() {
			res := r.db.Where("status = ?", status).Where("formed_at > ?", formDateStart).
				Find(&rocketFlights)
			if res.Error != nil {
				return rocketFlights, res.Error
			}

			rocketFlights, err = r.GetLoginsForFlights(rocketFlights)
			if err != nil {
				return rocketFlights, err
			}

			return rocketFlights, nil
		}

		// фильтрация по статусу, formDateStart и formDateEnd
		res := r.db.Where("status = ?", status).Where("formed_at BETWEEN ? AND ?", formDateStart, formDateEnd).Find(&rocketFlights)
		if res.Error != nil {
			return rocketFlights, res.Error
		}

		rocketFlights, err = r.GetLoginsForFlights(rocketFlights)
		if err != nil {
			return rocketFlights, err
		}

		return rocketFlights, nil
	}

	if formDateStart.IsZero() {
		if formDateEnd.IsZero() {
			// без фильтрации
			res := r.db.Where("status IN (?)", []string{"formed", "completed", "rejected"}).Find(&rocketFlights)
			if res.Error != nil {
				return rocketFlights, res.Error
			}

			rocketFlights, err = r.GetLoginsForFlights(rocketFlights)
			if err != nil {
				return rocketFlights, err
			}

			return rocketFlights, nil
		}

		// фильтрация по formDateEnd
		res := r.db.Where("status IN (?)", []string{"formed", "completed", "rejected"}).Where("formed_at < ?", formDateEnd).Find(&rocketFlights)
		if res.Error != nil {
			return rocketFlights, res.Error
		}

		rocketFlights, err = r.GetLoginsForFlights(rocketFlights)
		if err != nil {
			return rocketFlights, err
		}

		return rocketFlights, nil
	}

	if formDateEnd.IsZero() {
		// фильтрация по formDateStart
		res := r.db.Where("status IN (?)", []string{"formed", "completed", "rejected"}).Where("formed_at > ?", formDateStart).Find(&rocketFlights)
		if res.Error != nil {
			return rocketFlights, res.Error
		}

		rocketFlights, err = r.GetLoginsForFlights(rocketFlights)
		if err != nil {
			return rocketFlights, err
		}

		return rocketFlights, nil
	}

	//фильтрация по formDateStart и formDateEnd
	res := r.db.Where("status IN (?)", []string{"formed", "completed", "rejected"}).
		Where("formed_at BETWEEN ? AND ?", formDateStart, formDateEnd).Find(&rocketFlights)
	if res.Error != nil {
		return rocketFlights, res.Error
	}

	rocketFlights, err = r.GetLoginsForFlights(rocketFlights)
	if err != nil {
		return rocketFlights, err
	}

	return rocketFlights, nil
}

func (r *Repository) GetRocketFlightDraft(userId int) (int, error) {
	var rocketFlight models.RocketFlight
	err := r.db.First(&rocketFlight, "creator_id = ? and status = 'draft'", userId)
	if err.Error != nil && err.Error != gorm.ErrRecordNotFound {
		return 0, err.Error
	}

	return rocketFlight.FlightId, nil
}

func (r *Repository) GetRocketFlightById(flightId int) (models.RocketFlight, []models.Payload, error) {
	var rocketFlight models.RocketFlight
	// var rocketFlightDetailed models.RocketFlightDetailed
	var payloads []models.Payload
	var err error

	//информация по данному полёту
	result := r.db.First(&rocketFlight, "flight_id =?", flightId)
	if result.Error != nil {
		// log.Println("Ошибка при получении данного полёта")
		return models.RocketFlight{}, []models.Payload{}, result.Error
	}

	rocketFlights := []models.RocketFlight{rocketFlight}
	rocketFlights, err = r.GetLoginsForFlights(rocketFlights)
	if err != nil {
		return models.RocketFlight{}, []models.Payload{}, err
	}
	rocketFlight = rocketFlights[0]

	//полезные нагрузки КА, принятые на данный полёт
	result = r.db.Table("flights_payloads").Select("payloads.*").
		Joins("JOIN payloads ON flights_payloads.payload_id = payloads.payload_id").
		Where("flights_payloads.flight_id = ?", flightId).Find(&payloads)
	if result.Error != nil {
		// log.Println("Ошибка при получении заявок на полёт КА по данному полёту")
		return models.RocketFlight{}, []models.Payload{}, result.Error
	}

	return rocketFlight, payloads, nil
}

func (r *Repository) ChangeRocketFlight(changedRocketFlight models.RocketFlight) error {
	var oldRocketFlight models.RocketFlight
	result := r.db.Table("rocket_flights").Where("creator_id = ? AND status = ?", changedRocketFlight.CreatorId, "draft").Find(&oldRocketFlight)
	if result.Error != nil {
		return result.Error
	}

	if !changedRocketFlight.FlightDate.IsZero() {
		oldRocketFlight.FlightDate = changedRocketFlight.FlightDate
	}

	if changedRocketFlight.LoadCapacity != 0 {
		oldRocketFlight.LoadCapacity = changedRocketFlight.LoadCapacity
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
		return fmt.Errorf("Такой сформированной полезные нагрузки ракеты-носителя нет")
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
	if result.Error != nil && result.Error.Error() == "record not found" {
		return errors.New("Заявки-черновика на полёт ракеты-носителя нет")
	}
	if result.Error != nil {
		return result.Error
	}

	rocketFlight.Status = "deleted"
	result = r.db.Save(rocketFlight)
	return result.Error
}
