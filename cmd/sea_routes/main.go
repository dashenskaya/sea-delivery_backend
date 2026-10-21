package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"sea_routes/internal/api"
	"sea_routes/internal/app/dsn"
	"sea_routes/internal/app/repository"

	"context"
	"time"

	"sea_routes/internal/app/redisstore"
)

// @title Sea Routes API
// @version 1.0
// @description API сервиса морских маршрутов Sea Routes
// @description Регистрация, аутентификация, маршруты, публикация и лайки.

// @host localhost:8080
// @BasePath /api
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите Bearer и ваш JWT через пробел

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

	rep, err :=
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

	// Подключение к Redis.
	redisAddr := os.Getenv("REDIS_ADDR")

	if redisAddr == "" {
		log.Fatal("REDIS_ADDR не задан в .env")
	}

	redisStore := redisstore.New(redisAddr)
	defer redisStore.Close()

	// Проверяем соединение при запуске.
	redisCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := redisStore.Ping(redisCtx); err != nil {
		log.Fatal("Ошибка подключения к Redis: ", err)
	}

	log.Println("Redis connected successfully!")

	api.StartServer(rep, redisStore)

	log.Println("Application terminated!")
}
