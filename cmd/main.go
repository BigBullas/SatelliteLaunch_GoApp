package main

import (
	"RIP_lab1/internal/api/handler"
	"log"
)

func main() {
	log.Println("Application start!")

	h := handler.NewHandler(repo)
	h.StartServer()
	log.Println("Application terminated!")
}

//  docker run --name postgres -e POSTGRES_PASSWORD=password123 -d postgres
