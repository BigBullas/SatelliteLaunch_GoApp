package handler

import (
	"net/http"

	"RIP_lab1/internal/api"
	"RIP_lab1/internal/pkg/minio"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	repo   api.Repo
	minio  minio.Client
	logger *logrus.Entry
}

func NewHandler(repo api.Repo, minio minio.Client, logger *logrus.Logger) *Handler {
	return &Handler{repo: repo, minio: minio, logger: logger.WithField("component", "handler")}
}

func (h *Handler) StartServer() {
	h.logger.Println("Server start up")

	r := gin.Default()
	r.GET("/ping", h.Ping)

	r.LoadHTMLGlob("templates/*")

	r.Static("/image", "./resources")
	r.Static("/style", "./style")
	// полезные нагрузки
	r.GET("/payloads", h.GetPayloadList)
	r.GET("/payloads/:id", h.GetCardPayloadById)
	r.POST("/payloads", h.CreateNewPayload)
	r.PUT("/payloads/:id", h.ChangePayload)
	r.DELETE("/payloads/:id", h.DeletePayloadById)

	// удалить после перехода на фронт
	r.POST("/payload/:id", h.DeletePayloadById)

	// полёты ракет-носителей
	r.GET("/rocket_flights", h.GetRocketFlightList)
	r.GET("/rocket_flights/:id", h.GetRocketFlightById)
	r.PUT("/rocket_flights", h.ChangeRocketFlight)
	r.PUT("/rocket_flights/form", h.FormRocketFlight)
	r.PUT("/rocket_flights/:id/response", h.ResponceRocketFlight)
	r.DELETE("/rocket_flights", h.DeleteRocketFlight)

	// формирование информации о будущем полёте через полезные нагрузки
	r.POST("/payloads/rocket_flight", h.AddPayloadToFlight)

	// m-m
	r.DELETE("/flights_payloads/payload/:id", h.DeletePayloadFromFlight)
	r.PUT("/flights_payloads/payload/:id/count/:count", h.ChangeCountFlightsPayload)

	err := r.Run(":8080")
	if err != nil {
		h.logger.Fatalln(err)
	}
}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "pong",
		})
}
