package models

type FlightsPayload struct {
	FlightId        int	`json:"flight_id"`
	PayloadId       int	`json:"payload_id"`
	CountSatellites int	`json:"count_satellites"`
}
