package main

import (
	"RIP_lab1/internal/api/handler"
	"time"

	_ "RIP_lab1/docs"

	"github.com/sirupsen/logrus"
)

// @title Satellite launch application
// @version 1.0
// @description Web application on the topic "Launching satellites from the Vostochny cosmodrome"

// @host localhost:8080
// @schemes http
// @BasePath /
func main() {
	logger := logrus.New()
	formatter := &logrus.TextFormatter{
		TimestampFormat: time.DateTime,
		FullTimestamp:   true,
	}
	logger.SetFormatter(formatter)

	handler := handler.NewHandler(logger)
	r := handler.InitRoutes()
	err := r.Run("192.168.86.193:8080")
	if err != nil {
		logger.Fatalln(err)
	}
}
