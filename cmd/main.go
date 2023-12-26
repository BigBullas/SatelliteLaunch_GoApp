package main

import (
	"RIP_lab1/internal/api/handler"
	"RIP_lab1/internal/api/repository"
	"RIP_lab1/internal/pkg"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Println("Application start!")

	dsn, err := pkg.GetConnectionString()
	if err != nil {
		log.Error(err)
	}
	log.Info()

	repo, err := repository.NewRepo(dsn)
	if err != nil {
		log.Error(err)
	}

	h := handler.NewHandler(repo)
	h.StartServer()
	log.Println("Application terminated!")
}
