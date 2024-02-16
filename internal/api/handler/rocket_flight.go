package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"RIP_lab1/internal/models"
)

// GetRocketFlightList godoc
// @Summary Get rocket flight list
// @Description Retrieve a list of rocket flights based on the provided query parameters.
// @Tags RocketFlights
// @Accept json
// @Produce json
// @Param form_date_start query string false "Start date of the formation period"
// @Param form_date_end query string false "End date of the formation period"
// @Param status query string false "Status of the flight"
// @Success 200 {array} []models.RocketFlight
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /rocket_flights [get]
func (h *Handler) GetRocketFlightList(c *gin.Context) {
	queryString := c.Request.URL.Query()
	strFormDateStart := queryString.Get("form_date_start")
	strFormDateEnd := queryString.Get("form_date_end")
	strStatus := queryString.Get("status")

	var formDateStart, formDateEnd time.Time
	var err error

	if strFormDateStart != "" {
		formDateStart, err = time.Parse("2006-01-02", strFormDateStart)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Неверно указана дата конца периода формирования полёта"})
			return
		}
	}

	if strFormDateEnd != "" {
		formDateEnd, err = time.Parse("2006-01-02", strFormDateEnd)
		if err != nil {
			h.logger.Println("err_parsing: ", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Неверно указана дата конца периода формирования полёта"})
			return
		}

		if formDateEnd.Before(formDateStart) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Дата конца периода формирования полёта раньше даты начала"})
			return
		}
	}

	userId := c.GetInt(userCtx)
	isAdmin := c.GetBool(adminCtx)

	rocketFlights, err := h.repo.GetRocketFlightList(formDateStart, formDateEnd, strStatus, userId, isAdmin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	for i := range rocketFlights {

		h.logger.Println("in for: ", i)

		rocketFlights[i].ModeratorId = 0
		rocketFlights[i].CreatorId = 0
	}

	c.JSON(http.StatusOK, rocketFlights)
}

// GetRocketFlightById godoc
// @Summary Get rocket flight by ID
// @Description Retrieve a rocket flight and its associated payloads based on the provided ID.
// @Tags RocketFlights
// @Accept json
// @Produce json
// @Param id path int true "ID of the rocket flight"
// @Success 200 {object} []models.Payload
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /rocket_flights/{id} [get]
func (h *Handler) GetRocketFlightById(c *gin.Context) {
	strFlightId := c.Param("id")
	flightId, err := strconv.Atoi(strFlightId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Ошибка при преобразовании id полёта в число"})
		return
	}

	rocket_flight, payloads, err := h.repo.GetRocketFlightById(flightId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	rocket_flight.CreatorId = 0
	rocket_flight.ModeratorId = 0

	c.JSON(http.StatusOK, gin.H{"rocket_flight": rocket_flight, "payloads": payloads})
}

// ChangeRocketFlight godoc
// @Summary Change rocket flight
// @Description Update the details of a rocket flight.
// @Tags RocketFlights
// @Accept json
// @Produce json
// @Param flightDetails body models.RocketFlight true "Details of the rocket flight"
// @Success 200 {string} string "Rocket flight details successfully updated"
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /rocket_flights [put]
func (h *Handler) ChangeRocketFlight(c *gin.Context) {
	newRocketFlight := models.RocketFlight{}

	err := c.BindJSON(&newRocketFlight)

	h.logger.Println(newRocketFlight)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	newRocketFlight.CreatorId = c.GetInt(userCtx)

	h.logger.Println("RocketFlightData", newRocketFlight.FlightDate)

	if !(newRocketFlight.FlightDate == time.Time{}) {
		_, err = time.Parse("2006-01-02  15:04:05 -0700 MST", newRocketFlight.FlightDate.String())

		if err != nil {
			h.logger.Println("error: ", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

	}

	err = h.repo.ChangeRocketFlight(newRocketFlight)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Информация о полёте успешно изменена"})
}

// FormRocketFlight godoc
// @Summary Form rocket flight
// @Description Form a rocket flight.
// @Tags RocketFlights
// @Accept json
// @Produce json
// @Param flightStatus body models.RocketFlight true "Details of the rocket flight"
// @Success 200 {string} string "Rocket flight successfully formed"
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /rocket_flights/form [post]
func (h *Handler) FormRocketFlight(c *gin.Context) {
	var newFlightStatus models.RocketFlight
	var flightId int
	err := c.BindJSON(&newFlightStatus)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if newFlightStatus.Status != "formed" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Поменять статус можно только на 'formed'"})
		return
	}

	newFlightStatus.CreatorId = c.GetInt(userCtx)

	flightId, err = h.repo.FormRocketFlight(newFlightStatus)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	err = h.StartScanning(flightId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"mesage": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Статус успешно изменен на 'formed'"})
}

// ResponceRocketFlight godoc
// @Summary Response rocket flight
// @Description Update the status of a rocket flight.
// @Tags RocketFlights
// @Accept json
// @Produce json
// @Param id path int true "ID of the rocket flight"
// @Param flightStatus body models.RocketFlight true "New status of the rocket flight"
// @Success 200 {string} string "Rocket flight status successfully updated"
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /rocket_flights/{id}/response [put]
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
	newFlightStatus.ModeratorId = c.GetInt(userCtx)

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

// DeleteRocketFlight godoc
// @Summary Delete rocket flight
// @Description Delete a rocket flight draft.
// @Tags RocketFlights
// @Accept json
// @Produce json
// @Success 200 {string} string "Rocket flight draft successfully deleted"
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /rocket_flights [delete]
func (h *Handler) DeleteRocketFlight(c *gin.Context) {
	userId := c.GetInt(userCtx)
	err := h.repo.DeleteRocketFlight(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Заявка-черновик успешно удалена"})
}

func (h *Handler) StartScanning(flightId int) error {
	var flightAsync models.FlightAsync
	flightAsync.Id = flightId
	body, _ := json.Marshal(flightAsync)

	asyncServer := &http.Client{}
	req, err := http.NewRequest("PUT", "http://127.0.0.1:8081/calculate/", bytes.NewBuffer(body))
	if err != nil {
		return errors.New("Ошибка при создании запроса в бухгалтерию")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := asyncServer.Do(req)
	if err != nil {
		return errors.New("ошибка при отправке запроса в бухгалтерию")
	}

	if resp.StatusCode == 200 {
		return nil
	}

	return errors.New("Запрос не принят в обработку")
}

func (h *Handler) FinishCalculating(c *gin.Context) {
	const token = "qwertyuioplkjhgfdsa0987654321"
	var flightAsync models.FlightAsync
	if err := c.BindJSON(&flightAsync); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		h.logger.Println(err)
		return
	}
	if flightAsync.Token != token {
		c.AbortWithError(http.StatusForbidden, errors.New("Неверный токен"))
		return
	}
	err := h.repo.FinishCalculating(flightAsync)
	if err != nil {
		h.logger.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Данные сохранены"})
}
