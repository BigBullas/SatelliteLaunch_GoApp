package handler

import (
	"log"
	"net/http"
	"strconv"

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

	r.GET("/home", h.GetRequestForDeliveryList)
	r.GET("/card/", h.GetCardRequestForDeliveryById)
	r.POST("/card/:cardId", h.DeleteRequestForDeliveryById)

	r.Static("/image", "./resources")
	r.Static("/style", "./style")

	err := r.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func (h *Handler) GetRequestForDeliveryList(c *gin.Context) {
	queryString := c.Request.URL.Query()       // queryString - это тип url.Values, который содержит все query параметры
	strSearch := queryString.Get("spacecraft") // Получение значения конкретного параметра по его имени

	data, err := h.repo.GetRequestForDeliveryList(strSearch)
	if err != nil {
		log.Println(err)
	}

	c.HTML(http.StatusOK, "index.gohtml", gin.H{
		"cards":      data,
		"spacecraft": strSearch,
	})
}

func (h *Handler) GetCardRequestForDeliveryById(c *gin.Context) {
	queryString := c.Request.URL.Query() // queryString - это тип url.Values, который содержит все query параметры

	strCardId := queryString.Get("cardId") // Получение значения конкретного параметра по его имени
	cardId, err := strconv.Atoi(strCardId)
	if err != nil {
		log.Println("Ошибка при преобразовании строки в число:", err)
		return
	}

	card, err := h.repo.GetCardRequestForDeliveryByID(cardId)
	if err != nil {
		log.Println(err)
	}

	c.HTML(http.StatusOK, "card_launch_vehicle.gohtml", gin.H{
		"card": card,
	})
}

func (h *Handler) DeleteRequestForDeliveryById(c *gin.Context) {

	strCardId := c.Param("cardId")
	cardId, err := strconv.Atoi(strCardId)
	if err != nil {
		log.Println("Ошибка при преобразовании строки в число:", err)
	}

	log.Println("HANDLER, cardId: ", cardId)

	err = h.repo.DeleteRequestForDeliveryById(cardId)
	if err != nil {
		log.Println("Ошибка при получении заявки на доставку по id: ", cardId, err)
		c.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/home")
}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "pong",
		})
}
