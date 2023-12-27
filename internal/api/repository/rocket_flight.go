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
