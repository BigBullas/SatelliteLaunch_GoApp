package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"RIP_lab1/internal/models"
)

func (h *Handler) GetRocketFlightList(c *gin.Context) {
	queryString := c.Request.URL.Query()
	strFormDateStart := queryString.Get("form_date_start")
	strFormDateEnd := queryString.Get("form_date_end")
	strStatus := queryString.Get("status")

	var formDateStart, formDateEnd time.Time
	var err error

	if strFormDateStart != "" {
		formDateStart, err = time.Parse("2006-01-02 15:04:05", strFormDateStart)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Неверно указана дата конца периода формирования полёта"})
			return
		}
	}

	if strFormDateEnd != "" {
		formDateEnd, err = time.Parse("2006-01-02 15:04:05", strFormDateEnd)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Неверно указана дата конца периода формирования полёта"})
			return
		}

		if formDateEnd.Before(formDateStart) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Дата конца периода формирования полёта раньше даты начала"})
			return
		}
	}

	rocketFlights, err := h.repo.GetRocketFlightList(formDateStart, formDateEnd, strStatus)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, rocketFlights)
}

func (h *Handler) GetRocketFlightById(c *gin.Context) {
	strFlightId := c.Param("id")
	flightId, err := strconv.Atoi(strFlightId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Ошибка при преобразовании id полёта в число"})
		return
	}

	rocket_flight, flight_requests, err := h.repo.GetRocketFlightById(flightId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"rocket_flight": rocket_flight, "flight_requests": flight_requests})
}

func (h *Handler) ChangeRocketFlight(c *gin.Context) {
	var newRocketFlight models.RocketFlightChangeable
	err := c.BindJSON(&newRocketFlight)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	newRocketFlight.CreatorId = 1

	err = h.repo.ChangeRocketFlight(newRocketFlight)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Информация о полёте успешно изменена"})
}

func (h *Handler) FormRocketFlight(c *gin.Context) {
	var newFlightStatus models.RocketFlight
	err := c.BindJSON(&newFlightStatus)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if newFlightStatus.Status != "formed" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Поменять статус можно только на 'formed'"})
		return
	}

	newFlightStatus.CreatorId = 1

	err = h.repo.FormRocketFlight(newFlightStatus)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Статус успешно изменен на 'formed'"})
}

func (h *Handler) ResponceRocketFlight(c *gin.Context) {
	var newFlightStatus models.RocketFlight
	err := c.BindJSON(&newFlightStatus)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	strFlightId := c.Param("id")
	flightId, err := strconv.Atoi(strFlightId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	newFlightStatus.FlightId = flightId
	newFlightStatus.ModeratorId = 2

	if newFlightStatus.Status != "completed" && newFlightStatus.Status != "rejected" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Поменять статус можно только на 'completed или 'rejected'"})
		return
	}
	err = h.repo.ResponceRocketFlight(newFlightStatus)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Статус заявки успешно изменен"})
	return
}
