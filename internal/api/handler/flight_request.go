package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"RIP_lab1/internal/models"
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

func (h *Handler) CreateNewRequestForFlight(c *gin.Context) {
	var newFlightRequest models.FlightRequest

	// type FlightRequest struct {
	// 	Id              int `gorm:"primarykey"`
	// 	ImgURL          string
	// 	Title           string
	// 	LoadCapacity    float64
	// 	Description     string N
	// 	DetailedDesc    string N
	// 	DesiredPrice    float64 N
	// 	FlightDateStart time.Time
	// 	FlightDateEnd   time.Time
	// }

	newFlightRequest.Title = c.Request.FormValue("title")
	if newFlightRequest.Title == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Название КА не может быть пустым"})
		return
	}

	// log.Println("title", newFlightRequest.Title)

	_, header, err := c.Request.FormFile("image")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	if header == nil || header.Size == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Не было выслано изображение"})
		return
	}
	newFlightRequest.ImgURL = "https://ntv-static.cdnvideo.ru/home/news/2023/20230205/sputn_io.jpg"

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

	// if newFlightRequest.Image, err = h.minio.SaveImage(c.Request.Context(), file, header); err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "ошибка при сохранении изображения"})
	// 	return
	// }

	err = h.repo.CreateNewRequestForFlight(newFlightRequest)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusCreated, "Новая заявка на полёт успешно создана")
}

func (h *Handler) ChangeRequestForFlight(c *gin.Context) {
	var changedFlightRequest models.FlightRequest

	strCardId := c.Param("id")
	cardId, err := strconv.Atoi(strCardId)
	if err != nil {
		log.Println("Ошибка при преобразовании строки в число:", err)
		return
	}
	changedFlightRequest.RequestId = cardId

	// type FlightRequest struct {
	// 	Id              int `gorm:"primarykey"`
	// 	ImgURL          string
	// 	Title           string
	// 	LoadCapacity    float64
	// 	Description     string N
	// 	DetailedDesc    string N
	// 	DesiredPrice    float64 N
	// 	FlightDateStart time.Time
	// 	FlightDateEnd   time.Time
	// }

	changedFlightRequest.Title = c.Request.FormValue("title")

	// log.Println("title", newFlightRequest.Title)

	_, header, _ := c.Request.FormFile("image")
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
	// 	return
	// }
	if header != nil && header.Size != 0 {
		// delete old image
		// add new image
		changedFlightRequest.ImgURL = "https://finobzor.ru/uploads/posts/2016-09/org_vrke626.jpg"
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusCreated, "Заявка на полёт успешно изменена")

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
