package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetRequestForFlightList(c *gin.Context) {
	queryString := c.Request.URL.Query()            // queryString - это тип url.Values, который содержит все query параметры
	strSearch := queryString.Get("space_satellite") // Получение значения конкретного параметра по его имени

	data, err := h.repo.GetRequestForFlightList(strSearch)
	if err != nil {
		log.Println(err)
	}

	c.HTML(http.StatusOK, "index.gohtml", gin.H{
		"cards":           data,
		"space_satellite": strSearch,
	})
}

func (h *Handler) GetCardRequestForFlightById(c *gin.Context) {
	// queryString := c.Request.URL.Query() // queryString - это тип url.Values, который содержит все query параметры

	// strCardId := queryString.Get("id") // Получение значения конкретного параметра по его имени

	strCardId := c.Param("id")
	cardId, err := strconv.Atoi(strCardId)
	if err != nil {
		log.Println("Ошибка при преобразовании строки в число:", err)
		return
	}

	card, err := h.repo.GetCardRequestForFlightById(cardId)
	if err != nil {
		log.Println(err)
	}

	c.HTML(http.StatusOK, "card_launch_vehicle.gohtml", gin.H{
		"card": card,
	})
}

func (h *Handler) DeleteRequestForFlightById(c *gin.Context) {

	strCardId := c.Param("id")
	cardId, err := strconv.Atoi(strCardId)
	if err != nil {
		log.Println("Ошибка при преобразовании строки в число:", err)
	}

	log.Println("HANDLER, id: ", cardId)

	err = h.repo.DeleteRequestForFlightById(cardId)
	if err != nil {
		log.Println("Ошибка при получении заявки на доставку по id: ", cardId, err)
		c.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/home")
}
