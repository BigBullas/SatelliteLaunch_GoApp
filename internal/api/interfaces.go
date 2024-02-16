package api

import (
	"RIP_lab1/internal/models"
	"context"
	"time"
)

type Repo interface {
	GetPayloadList(string, string, string, string, string) ([]models.Payload, error)
	GetCardPayloadById(int) (models.Payload, error)
	CreateNewPayload(models.Payload) error
	ChangePayload(models.Payload) error
	DeletePayloadById(int) error
	GetPayloadImageUrl(int) string

	AddPayloadToFlight(int, int) (bool, error)
	DeletePayloadFromFlight(int, int) error
	ChangeCountFlightsPayload(int, int, int) error

	GetRocketFlightList(time.Time, time.Time, string, int, bool) ([]models.RocketFlight, error)
	GetRocketFlightDraft(int) (int, error)
	GetRocketFlightById(int) (models.RocketFlight, []models.Payload, error)
	ChangeRocketFlight(models.RocketFlight) error
	FormRocketFlight(models.RocketFlight) (int, error)
	ResponceRocketFlight(models.RocketFlight) error
	DeleteRocketFlight(int) error
	FinishCalculating(models.FlightAsync) error

	SignUp(ctx context.Context, newUser models.User) (int, bool, error)
	GetByCredentials(context.Context, models.User) (models.User, error)
	GetUserInfo(context.Context, models.User) (models.User, error)
	ChangeProfile(context.Context, models.User) error
}
