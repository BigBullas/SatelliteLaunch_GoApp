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
	r.GET("/home", h.GetRequestForFlightList)
	r.GET("/flight_request/:id", h.GetCardRequestForFlightById)
	r.POST("/flight_request", h.CreateNewRequestForFlight)
	r.PUT("/flight_request/:id", h.ChangeRequestForFlight)
	r.DELETE("/flight_request/:id", h.DeleteRequestForFlightById)

	// удалить после перехода на фронт
	r.POST("/flight_request/:id", h.DeleteRequestForFlightById)

	// полёты ракет-носителей
	r.GET("/rocket_flights", h.GetRocketFlightList)
	r.GET("/rocket_flight/:id", h.GetRocketFlightById)
	r.PUT("/rocket_flight", h.ChangeRocketFlight)
	r.PUT("/rocket_flight/form", h.FormRocketFlight)
	r.PUT("/rocket_flight/:id/response", h.ResponceRocketFlight)
	r.DELETE("/rocket_flights", h.DeleteRocketFlight)

	// формирование информации о будущем полёте через заявки на полёт
	r.POST("/flight_request/rocket_flight", h.AddFlightRequestToFlight)

	// m-m
	r.DELETE("/flights_flight_requests/flight_request/:id", h.DeleteRequestFromFlight)

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
