package models

import "time"

type RequestForDelivery struct {
	ID                  uint `gorm:"primarykey"`
	ImgURL              string
	Title               string
	LoadCapacity        float64
	Description         string
	DetailedDescription string
	DesiredPrice        float64
	FlightDateStart     time.Time
	FlightDateEnd       time.Time
}

// request_id          integer not null
// constraint request_pk primary key,
// is_available        boolean,
// img_url             TEXT,
// title               varchar(100),
// load_capacity       float,
// description         TEXT,
// detailed_desc       TEXT,
// desired_price       float,
// flight_date_start   timestamp,
// flight_date_end     timestamp
