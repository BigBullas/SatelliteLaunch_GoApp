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

func (h *Handler) GetPayloadList(c *gin.Context) {
	queryString := c.Request.URL.Query()            // queryString - это тип url.Values, который содержит все query параметры
	strSearch := queryString.Get("space_satellite") // Получение значения конкретного параметра по его имени

	data, err := h.repo.GetPayloadList(strSearch)
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
		c.JSON(http.StatusOK, gin.H{"payloads": data, "draftRocketFlightId": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payloads": data, "draftRocketFlightId": flightId})
}

func (h *Handler) GetCardPayloadById(c *gin.Context) {
	strCardId := c.Param("id")
	cardId, err := strconv.Atoi(strCardId)
	if err != nil {
		log.Println("Ошибка при преобразовании строки в число:", err)
		return
	}

	card, err := h.repo.GetCardPayloadById(cardId)
	if err != nil {
		log.Println(err)
	}

	c.HTML(http.StatusOK, "card_launch_vehicle.gohtml", gin.H{
		"card": card,
	})
}

func (h *Handler) CreateNewPayload(c *gin.Context) {
	var newPayload models.Payload

	newPayload.Title = c.Request.FormValue("title")
	if newPayload.Title == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Название КА не может быть пустым"})
		return
	}

	// log.Println("title", newPayload.Title)

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	// if header == nil || header.Size == 0 {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Не было выслано изображение"})
	// 	return
	// }
	// newPayload.ImgURL = "https://ntv-static.cdnvideo.ru/home/news/2023/20230205/sputn_io.jpg"

	// log.Println("image", newPayload.ImgURL)

	loadCapacity := c.Request.FormValue("load_capacity")
	if loadCapacity != "" {
		newPayload.LoadCapacity, err = strconv.ParseFloat(loadCapacity, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Неверно указан полезный вес КА"})
			return
		}
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Полезный вес КА не может быть пустым"})
		return
	}

	// log.Println("load_capacity", newPayload.LoadCapacity)

	newPayload.Description = c.Request.FormValue("description")
	newPayload.DetailedDesc = c.Request.FormValue("detailed_description")

	// log.Println("descriptions: ", newPayload.Description, newPayload.DetailedDesc)

	desiredPrice := c.Request.FormValue("desired_price")
	if desiredPrice != "" {
		newPayload.DesiredPrice, err = strconv.ParseFloat(desiredPrice, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Неверно указана желаемая цена услуги"})
			return
		}
	}

	// log.Println("desired price: ", newPayload.DesiredPrice)

	startDate := c.Request.FormValue("flight_date_start")
	if startDate == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Дата начала желаемого периода полёта не может быть пустой"})
		return
	}

	// log.Println("start date: ", startDate)

	newPayload.FlightDateStart, err = time.Parse("2006-01-02 15:04:05", startDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Неверно указана дата начала желаемого периода полёта"})
		return
	}

	endDate := c.Request.FormValue("flight_date_end")
	if endDate == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Дата конца желаемого периода полёта не может быть пустой"})
		return
	}
	newPayload.FlightDateEnd, err = time.Parse("2006-01-02 15:04:05", endDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Неверно указана дата конца желаемого периода полёта"})
		return
	}

	newPayload.IsAvailable = true

	newPayload.ImgURL, err = h.minio.SaveImage(c.Request.Context(), file, header)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Ошибка при сохранении изображения"})
		return
	}

	err = h.repo.CreateNewPayload(newPayload)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusCreated, "Новая полезная нагрузка успешно создана")
}

func (h *Handler) ChangePayload(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		h.logger.Errorf("Handler/payload/ChangePayload/Ошибка при чтении файла: %s", err)
	}
	var changedPayload models.Payload

	strCardId := c.Param("id")
	cardId, err := strconv.Atoi(strCardId)
	if err != nil {
		log.Println("Ошибка при преобразовании строки в число:", err)
		return
	}
	changedPayload.PayloadId = cardId

	changedPayload.Title = c.Request.FormValue("title")

	// log.Println("title", newPayload.Title)

	if header != nil && header.Size != 0 {
		changedPayload.ImgURL, err = h.minio.SaveImage(c.Request.Context(), file, header)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}

		url := h.repo.GetPayloadImageUrl(changedPayload.PayloadId)

		// delete image from minio
		h.minio.DeleteImage(c.Request.Context(), utils.ExtractObjectNameFromUrl(url))
	}

	// log.Println("image", newPayload.ImgURL)

	loadCapacity := c.Request.FormValue("load_capacity")
	if loadCapacity != "" {
		changedPayload.LoadCapacity, err = strconv.ParseFloat(loadCapacity, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Неверно указан полезный вес КА"})
			return
		}
	}

	// log.Println("load_capacity", newPayload.LoadCapacity)

	changedPayload.Description = c.Request.FormValue("description")
	changedPayload.DetailedDesc = c.Request.FormValue("detailed_description")

	// log.Println("descriptions: ", newPayload.Description, newPayload.DetailedDesc)

	desiredPrice := c.Request.FormValue("desired_price")
	if desiredPrice != "" {
		changedPayload.DesiredPrice, err = strconv.ParseFloat(desiredPrice, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Неверно указана желаемая цена услуги"})
			return
		}
	}

	// log.Println("desired price: ", newPayload.DesiredPrice)

	startDate := c.Request.FormValue("flight_date_start")
	if startDate != "" {
		changedPayload.FlightDateStart, err = time.Parse("2006-01-02 15:04:05", startDate)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Неверно указана дата начала желаемого периода полёта"})
			return
		}
	}

	// log.Println("start date: ", startDate)

	endDate := c.Request.FormValue("flight_date_end")
	if endDate != "" {
		changedPayload.FlightDateEnd, err = time.Parse("2006-01-02 15:04:05", endDate)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Неверно указана дата конца желаемого периода полёта"})
			return
		}
	}

	err = h.repo.ChangePayload(changedPayload)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "Полезная нагрузка успешно изменена")

}

func (h *Handler) DeletePayloadById(c *gin.Context) {
	strCardId := c.Param("id")
	cardId, err := strconv.Atoi(strCardId)

	h.logger.Print("strCardI: ", strCardId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = h.repo.DeletePayloadById(cardId)
	if err != nil {
		log.Println("Ошибка при получении заявки на доставку по id: ", cardId, err)
		c.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/payloads")
}

func (h *Handler) AddPayloadToFlight(c *gin.Context) {
	var creatorId int
	var payloadId int

	type RocketFlightShort struct {
		CreatorId int
		PayloadId int
	}

	jsonStr := RocketFlightShort{}

	err := c.ShouldBindJSON(&jsonStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	creatorId = 1
	payloadId = jsonStr.PayloadId

	if payloadId == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Требуется хотя бы одна полезная нагрузка"})
		return
	}

	err = h.repo.AddPayloadToFlight(creatorId, payloadId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Полезная нагрузка добавлена в планируемый полёт"})
	return
}

func (h *Handler) DeletePayloadFromFlight(c *gin.Context) {
	strPayloadId := c.Param("id")
	payloadId, err := strconv.Atoi(strPayloadId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	userId := 1

	err = h.repo.DeletePayloadFromFlight(userId, payloadId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Полезная нагрузка КА успешно удалена из планируемого полёта"})
}

func (h *Handler) ChangeCountFlightsPayload(c *gin.Context) {
	strPayloadId := c.Param("id")
	payloadId, err := strconv.Atoi(strPayloadId)
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

	err = h.repo.ChangeCountFlightsPayload(userId, payloadId, count)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Количество полезных нагрузок успешно изменено"})
}
