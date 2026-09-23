package main

import (
	"log"

	"github.com/joho/godotenv"

	"sea_routes/internal/api"
	"sea_routes/internal/app/dsn"
	"sea_routes/internal/app/repository"
)

func main() {
	log.Println("Application start!")

	err := godotenv.Load()

	if err != nil {
		log.Fatal(
			"Не удалось загрузить .env: ",
			err,
		)
	}

	rep, err :=
		repository.New(
			dsn.FromEnv(),
		)

	if err != nil {
		log.Fatal(
			"Не удалось подключиться к БД: ",
			err,
		)
	}

	api.StartServer(rep)

	log.Println("Application terminated!")
}
