package api

import (
	"log"

	"github.com/gin-gonic/gin"

	"sea_delivery/internal/app/handler"
)

func StartServer() {
	log.Println("Server start up")

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")

	router.GET("/sea-routes/feed", handler.SeaRouteFeedHandler)
	router.GET("/sea-routes/add", handler.SeaRouteAddHandler)
	router.GET("/sea-routes/catalog", handler.SeaRouteCatalogHandler)

	err := router.Run()

	if err != nil {
		log.Println(err)
	}

	log.Println("Server down")
}
