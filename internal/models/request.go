package models

import (
	"time"
)

type FlightRequest struct {
	RequestId       int `gorm:"primarykey"`
	ImgURL          string
	Title           string
	LoadCapacity    float64
	Description     string
	DetailedDesc    string
	DesiredPrice    float64
	FlightDateStart time.Time
	FlightDateEnd   time.Time
}
