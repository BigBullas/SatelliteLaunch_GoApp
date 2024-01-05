package main

import (
	"RIP_lab1/internal/api/handler"
	"RIP_lab1/internal/api/repository"
	"RIP_lab1/internal/pkg"
	"RIP_lab1/internal/pkg/minio"
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logger := logrus.New()
	formatter := &logrus.TextFormatter{
		TimestampFormat: time.DateTime,
		FullTimestamp:   true,
	}
	logger.SetFormatter(formatter)

	logger.Println("Application start!")

	dsn, err := pkg.GetConnectionString()
	if err != nil {
		logger.Error(err)
	}
	logger.Info(dsn)

	vp := viper.New()
	if err := initConfig(vp); err != nil {
		logger.Fatalf("main.go/Error initializing configs: %s", err.Error())
	}

	minioConfig := minio.InitConfig(vp)

	minioClient, err := minio.NewMinioClient(context.Background(), minioConfig, logger)
	if err != nil {
		logger.Fatalf("main.go/minioClientCreation: %s", err.Error())
	}

	repo, err := repository.NewRepo(dsn)
	if err != nil {
		logger.Fatalf("main.go/newRepoCreation: %s", err.Error())
	}

	h := handler.NewHandler(repo, minioClient, logger)
	h.StartServer()
	logger.Println("Application terminated!")
}

func initConfig(vp *viper.Viper) error {
	vp.AddConfigPath("./config")
	vp.SetConfigName("config")

	return vp.ReadInConfig()
}
