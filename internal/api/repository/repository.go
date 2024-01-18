package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"RIP_lab1/internal/models"
)

type Repository struct {
	db *gorm.DB
	logger *logrus.Entry
}

func NewRepo(logger *logrus.Logger, vp *viper.Viper) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(vp.GetString("db.connection_string")), &gorm.Config{})
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	err = db.AutoMigrate(&models.Payload{})
	err = db.AutoMigrate(&models.RocketFlight{})
	err = db.AutoMigrate(&models.FlightsPayload{})
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		panic("Миграция БД не удалась")
	}

	return &Repository{
		db: db,
		logger: logger.WithField("component", "repository"),
	}, nil
}
