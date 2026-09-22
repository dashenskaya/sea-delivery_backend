package api

import (
	"log"

	"github.com/gin-gonic/gin"

	"sea_routes/internal/app/handler"
)

func StartServer() {
	log.Println("Server start up")

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")

	router.GET("/sea_routes/feed", handler.SeaRouteFeedHandler)
	router.GET("/sea_routes/add", handler.SeaRouteAddHandler)
	router.GET("/sea_routes/catalog", handler.SeaRouteCatalogHandler)

	err := router.Run()

	if err != nil {
		log.Println(err)
	}

	log.Println("Server down")
}
