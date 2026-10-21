package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

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

	minioUseSSL, err :=
		strconv.ParseBool(
			os.Getenv("MINIO_USE_SSL"),
		)

	if err != nil {
		log.Fatal(
			"Некорректное значение MINIO_USE_SSL: ",
			err,
		)
	}

	_, err =
		repository.New(
			repository.RepositorySettings{
				PostgresDSN: dsn.FromEnv(),

				MinioEndpoint: os.Getenv(
					"MINIO_ENDPOINT",
				),

				MinioAccessKey: os.Getenv(
					"MINIO_ACCESS_KEY",
				),

				MinioSecretKey: os.Getenv(
					"MINIO_SECRET_KEY",
				),

				MinioBucket: os.Getenv(
					"MINIO_BUCKET",
				),

				MinioUseSSL: minioUseSSL,
			},
		)

	if err != nil {
		log.Fatal(
			"Не удалось инициализировать Repository: ",
			err,
		)
	}

	log.Println("Application terminated!")
}
