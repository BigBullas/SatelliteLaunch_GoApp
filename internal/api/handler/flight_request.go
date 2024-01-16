package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"RIP_lab1/internal/models"
	"RIP_lab1/internal/utils"
)

func (h *Handler) GetRequestForFlightList(c *gin.Context) {
	queryString := c.Request.URL.Query()            // queryString - это тип url.Values, который содержит все query параметры
	strSearch := queryString.Get("space_satellite") // Получение значения конкретного параметра по его имени

	data, err := h.repo.GetRequestForFlightList(strSearch)
	if err != nil {
		log.Println(err)
	}

	flightId, err := h.repo.GetRocketFlightDraft(1)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// if flightId == 0 {
	// 	c.HTML(http.StatusOK, "index.gohtml", gin.H{
	// 		"cards":           data,
	// 		"space_satellite": strSearch,
	// 	})
	// 	return
	// }

	// c.HTML(http.StatusOK, "index.gohtml", gin.H{
	// 	"cards":               data,
	// 	"space_satellite":     strSearch,
	// 	"draftRocketFlightId": flightId,
	// })

	if flightId == 0 {
		c.JSON(http.StatusOK, gin.H{"flight_requests": data, "draftRocketFlightId": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"flight_requests": data, "draftRocketFlightId": flightId})
}

func (h *Handler) GetCardRequestForFlightById(c *gin.Context) {
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

func (h *Handler) CreateNewRequestForFlight(c *gin.Context) {
	var newFlightRequest models.Payload

	newFlightRequest.Title = c.Request.FormValue("title")
	if newFlightRequest.Title == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Название КА не может быть пустым"})
		return
	}

	// log.Println("title", newFlightRequest.Title)

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// if header == nil || header.Size == 0 {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Не было выслано изображение"})
	// 	return
	// }
	// newFlightRequest.ImgURL = "https://ntv-static.cdnvideo.ru/home/news/2023/20230205/sputn_io.jpg"

	// log.Println("image", newFlightRequest.ImgURL)

	loadCapacity := c.Request.FormValue("load_capacity")
	if loadCapacity != "" {
		newFlightRequest.LoadCapacity, err = strconv.ParseFloat(loadCapacity, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Неверно указан полезный вес КА"})
			return
		}
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Полезный вес КА не может быть пустым"})
		return
	}

	// log.Println("load_capacity", newFlightRequest.LoadCapacity)

	newFlightRequest.Description = c.Request.FormValue("description")
	newFlightRequest.DetailedDesc = c.Request.FormValue("detailed_description")

	// log.Println("descriptions: ", newFlightRequest.Description, newFlightRequest.DetailedDesc)

	desiredPrice := c.Request.FormValue("desired_price")
	if desiredPrice != "" {
		newFlightRequest.DesiredPrice, err = strconv.ParseFloat(desiredPrice, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Неверно указана желаемая цена услуги"})
			return
		}
	}

	// log.Println("desired price: ", newFlightRequest.DesiredPrice)

	startDate := c.Request.FormValue("flight_date_start")
	if startDate == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Дата начала желаемого периода полёта не может быть пустой"})
		return
	}

	// log.Println("start date: ", startDate)

	newFlightRequest.FlightDateStart, err = time.Parse("2006-01-02 15:04:05", startDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Неверно указана дата начала желаемого периода полёта"})
		return
	}

	endDate := c.Request.FormValue("flight_date_end")
	if endDate == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Дата конца желаемого периода полёта не может быть пустой"})
		return
	}
	newFlightRequest.FlightDateEnd, err = time.Parse("2006-01-02 15:04:05", endDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Неверно указана дата конца желаемого периода полёта"})
		return
	}

	newFlightRequest.IsAvailable = true

	newFlightRequest.ImgURL, err = h.minio.SaveImage(c.Request.Context(), file, header)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "ошибка при сохранении изображения"})
		return
	}

	err = h.repo.CreateNewRequestForFlight(newFlightRequest)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusCreated, "Новая заявка на полёт успешно создана")
}

func (h *Handler) ChangeRequestForFlight(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		h.logger.Errorf("Handler/flight_request/ChangeRequestForFlight/Error read file: %s", err)
	}
	var changedFlightRequest models.Payload

	strCardId := c.Param("id")
	cardId, err := strconv.Atoi(strCardId)
	if err != nil {
		log.Println("Ошибка при преобразовании строки в число:", err)
		return
	}
	changedFlightRequest.PayloadId = cardId

	changedFlightRequest.Title = c.Request.FormValue("title")

	// log.Println("title", newFlightRequest.Title)

	if header != nil && header.Size != 0 {
		changedFlightRequest.ImgURL, err = h.minio.SaveImage(c.Request.Context(), file, header)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}

		url := h.repo.GetFlightRequestImageUrl(changedFlightRequest.PayloadId)

		// delete image from minio
		h.minio.DeleteImage(c.Request.Context(), utils.ExtractObjectNameFromUrl(url))
	}

	// log.Println("image", newFlightRequest.ImgURL)

	loadCapacity := c.Request.FormValue("load_capacity")
	if loadCapacity != "" {
		changedFlightRequest.LoadCapacity, err = strconv.ParseFloat(loadCapacity, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Неверно указан полезный вес КА"})
			return
		}
	}

	// log.Println("load_capacity", newFlightRequest.LoadCapacity)

	changedFlightRequest.Description = c.Request.FormValue("description")
	changedFlightRequest.DetailedDesc = c.Request.FormValue("detailed_description")

	// log.Println("descriptions: ", newFlightRequest.Description, newFlightRequest.DetailedDesc)

	desiredPrice := c.Request.FormValue("desired_price")
	if desiredPrice != "" {
		changedFlightRequest.DesiredPrice, err = strconv.ParseFloat(desiredPrice, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Неверно указана желаемая цена услуги"})
			return
		}
	}

	// log.Println("desired price: ", newFlightRequest.DesiredPrice)

	startDate := c.Request.FormValue("flight_date_start")
	if startDate != "" {
		changedFlightRequest.FlightDateStart, err = time.Parse("2006-01-02 15:04:05", startDate)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Неверно указана дата начала желаемого периода полёта"})
			return
		}
	}

	// log.Println("start date: ", startDate)

	endDate := c.Request.FormValue("flight_date_end")
	if endDate != "" {
		changedFlightRequest.FlightDateEnd, err = time.Parse("2006-01-02 15:04:05", endDate)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Неверно указана дата конца желаемого периода полёта"})
			return
		}
	}

	err = h.repo.ChangeRequestForFlight(changedFlightRequest)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "Заявка на полёт успешно изменена")

}

func (h *Handler) DeleteRequestForFlightById(c *gin.Context) {
	strCardId := c.Param("id")
	cardId, err := strconv.Atoi(strCardId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = h.repo.DeleteRequestForFlightById(cardId)
	if err != nil {
		log.Println("Ошибка при получении заявки на доставку по id: ", cardId, err)
		c.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/home")
}

func (h *Handler) AddFlightRequestToFlight(c *gin.Context) {
	var creatorId int
	var requestId int

	type RocketFlightShort struct {
		CreatorId int
		RequestId int
	}

	jsonStr := RocketFlightShort{}

	err := c.ShouldBindJSON(&jsonStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	creatorId = 1
	requestId = jsonStr.RequestId

	if requestId == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Требуется хотя бы одна заявка на полёт КА"})
		return
	}

	err = h.repo.AddFlightRequestToFlight(creatorId, requestId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Заявка на полёт КА добавлена в планируемый полёт"})
	return
}

func (h *Handler) DeleteRequestFromFlight(c *gin.Context) {
	strRequestId := c.Param("id")
	requestId, err := strconv.Atoi(strRequestId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	userId := 1

	err = h.repo.DeleteRequestFromFlight(userId, requestId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Заявка на полёт КА успешно удалена из планируемого полёта"})
}

func (h *Handler) ChangeCountFlightsFlightRequest(c *gin.Context) {
	strRequestId := c.Param("id")
	requestId, err := strconv.Atoi(strRequestId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	strCount := c.Param("count")
	count, err := strconv.Atoi(strCount)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	userId := 1

	err = h.repo.ChangeCountFlightsFlightRequest(userId, requestId, count)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Количество заявок на полёт КА успешно изменено"})
}
