package handler

import (
	"log"
	"net/http"

	"RIP_lab1/internal/api"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo api.Repo
}

func NewHandler(repo api.Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) StartServer() {
	log.Println("Server start up")

	r := gin.Default()
	r.GET("/ping", h.Ping)

	r.LoadHTMLGlob("templates/*")

	r.Static("/image", "./resources")
	r.Static("/style", "./style")
	// заявки на полёт
	r.GET("/flight_requests", h.GetRequestForFlightList)
	r.GET("/flight_requests/:id", h.GetCardRequestForFlightById)
	r.POST("/flight_requests", h.CreateNewRequestForFlight)
	r.PUT("/flight_requests/:id", h.ChangeRequestForFlight)
	r.DELETE("/flight_requests/:id", h.DeleteRequestForFlightById)

	// удалить после перехода на фронт
	r.POST("/flight_request/:id", h.DeleteRequestForFlightById)

	// полёты ракет-носителей
	r.GET("/rocket_flights", h.GetRocketFlightList)
	r.GET("/rocket_flights/:id", h.GetRocketFlightById)
	r.PUT("/rocket_flights", h.ChangeRocketFlight)
	r.PUT("/rocket_flights/form", h.FormRocketFlight)
	r.PUT("/rocket_flights/:id/response", h.ResponceRocketFlight)
	r.DELETE("/rocket_flights", h.DeleteRocketFlight)

	// формирование информации о будущем полёте через заявки на полёт
	r.POST("/flight_requests/rocket_flight", h.AddFlightRequestToFlight)

	// m-m
	r.DELETE("/flights_flight_requests/flight_request/:id", h.DeleteRequestFromFlight)
	r.PUT("/flights_flight_requests/flight_request/:id/count/:count", h.ChangeCountFlightsFlightRequest)

	err := r.Run(":8080")
	if err != nil {
		log.Fatalln(err)
	}
}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "pong",
		})
}
