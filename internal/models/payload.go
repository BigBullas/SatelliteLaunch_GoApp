package models

import (
	"time"
)

type Payload struct {
	PayloadId       int `gorm:"primarykey"`
	IsAvailable     bool
	ImgURL          string
	Title           string
	LoadCapacity    float64
	Description     string
	DetailedDesc    string
	DesiredPrice    float64
	FlightDateStart time.Time
	FlightDateEnd   time.Time
}
