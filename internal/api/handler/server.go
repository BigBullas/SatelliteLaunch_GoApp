package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"RIP_lab1/internal/api"
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
	r.GET("/card/:id", h.GetRequestForDeliveryList)
	r.POST("/card/:id", h.DeleteRequestForDeliveryById)

	r.Static("/image", "./resources")
	r.Static("/styles", "./styles")

	err := r.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func (h *Handler) GetRequestForDeliveryList(c *gin.Context) {
	queryString := c.Request.URL.Query() // queryString - это тип url.Values, который содержит все query параметры
	strSearch := queryString.Get("spacecraft") // Получение значения конкретного параметра по его имени

	data, err := h.repo.GetRequestForDeliveryList(strSearch)
	if err != nil {
		log.Println(err)
	}

	c.HTML(http.StatusOK, "main_page.gohtml", gin.H{
		"cards": data,
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

	c.HTML(http.StatusOK, "spacecraft_card.gohtml", gin.H{
		"card": card,
	})
}

func (h *Handler) DeleteRequestForDeliveryById(c *gin.Context) {
	queryString := c.Request.URL.Query() // queryString - это тип url.Values, который содержит все query параметры

	strCardId := queryString.Get("cardId") // Получение значения конкретного параметра по его имени
	cardId, err := strconv.Atoi(strCardId)
	if err != nil {
		log.Println("Ошибка при преобразовании строки в число:", err)
		return
	}
	
	err = h.repo.DeleteRequestForDeliveryById(cardId)
	if err != nil {
		log.Printf("Ошибка при получении заявки на доставку по id: ", cardId, err)
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