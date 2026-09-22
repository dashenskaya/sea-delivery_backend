package main

import (
	"log"

	"sea_routes/internal/api"
)

func main() {
	log.Println("Application start!")

	api.StartServer()

	log.Println("Application terminated!")
}
