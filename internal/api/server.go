package api

import (
	"log"

	"github.com/gin-gonic/gin"

	"sea_routes/internal/app/handler"
	"sea_routes/internal/app/repository"
)

func StartServer(rep *repository.Repository) {
	log.Println("Server start up")

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.Static(
		"/static",
		"./resources",
	)

	hand := handler.NewHandler(rep)

	router.GET(
		"/sea_routes/feed",
		hand.SeaRouteFeedHandler,
	)

	router.GET(
		"/sea_routes/add",
		hand.SeaRouteAddHandler,
	)

	router.GET(
		"/sea_routes/catalog",
		hand.SeaRouteCatalogHandler,
	)

	router.POST(
		"/sea_routes/create",
		hand.SeaRouteCreateHandler,
	)

	router.POST(
		"/sea_routes/publish",
		hand.SeaRoutePublishHandler,
	)

	router.POST(
		"/sea_routes/delete",
		hand.SeaRouteDeleteHandler,
	)

	err := router.Run()

	if err != nil {
		log.Println(err)
	}

	log.Println("Server down")
}
