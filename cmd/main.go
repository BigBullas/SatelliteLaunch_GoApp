package main

import (
	"RIP_lab1/internal/api"
	"log"
)

func main() {
	log.Println("Application start!")
	api.StartServer()
	log.Println("Application terminated!")
}

//  docker run --name postgres -e POSTGRES_PASSWORD=password123 -d postgres
