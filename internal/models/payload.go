package models

import (
	"time"
)

type Payload struct {
	PayloadId       int       `gorm:"primarykey" json:"payload_id"`
	IsAvailable     bool      `json:"is_available"`
	ImgURL          string    `json:"img_url"`
	Title           string    `json:"title"`
	LoadCapacity    float64   `json:"load_capacity"`
	Description     string    `json:"description"`
	DetailedDesc    string    `json:"detailed_desc"`
	DesiredPrice    float64   `json:"desired_price"`
	FlightDateStart time.Time `json:"flight_date_start"`
	FlightDateEnd   time.Time `json:"flight_date_end"`
	CountSatellites int       `json:"count,omitempty"`
}
