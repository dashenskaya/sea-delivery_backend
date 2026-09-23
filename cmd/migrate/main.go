package main

import (
	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"sea_routes/internal/app/ds"
	"sea_routes/internal/app/dsn"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Не удалось загрузить .env: ", err)
	}

	db, err := gorm.Open(
		postgres.Open(dsn.FromEnv()),
		&gorm.Config{},
	)

	if err != nil {
		log.Fatal("Не удалось подключиться к БД: ", err)
	}

	err = db.AutoMigrate(
		&ds.SeaRouteUser{},
		&ds.SeaRoute{},
		&ds.SeaRouteLike{},
	)

	if err != nil {
		log.Fatal("Ошибка миграции: ", err)
	}

	err = db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_one_draft_per_user
		ON sea_routes (sea_route_creator_id)
		WHERE sea_route_status = 'черновик'
	`).Error

	if err != nil {
		log.Fatal("Ошибка создания ограничения для черновика: ", err)
	}

	log.Println("Миграция выполнена успешно")
}
