package main

import (
	"sea_routes/internal/app/ds"
	"sea_routes/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()

	db, err :=
		gorm.Open(
			postgres.Open(dsn.FromEnv()),
			&gorm.Config{},
		)

	if err != nil {
		panic("failed to connect database")
	}

	err =
		db.AutoMigrate(
			&ds.SeaRouteUser{},
			&ds.SeaRoute{},
			&ds.SeaRouteLike{},
		)

	if err != nil {
		panic("cant migrate db")
	}

	// Один пользователь может иметь
	// только один черновик.
	err =
		db.Exec(`
			CREATE UNIQUE INDEX IF NOT EXISTS idx_one_draft_per_user
			ON sea_routes (sea_route_creator_id)
			WHERE sea_route_status = 'черновик'
		`).Error

	if err != nil {
		panic("cant create draft index")
	}
}
