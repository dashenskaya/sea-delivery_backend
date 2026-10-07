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

	apiGroup :=
		router.Group("/api")

	apiGroup.GET(
		"/sea-routes",
		hand.GetSeaRoutesAPI,
	)

	apiGroup.POST(
		"/sea-routes",
		hand.CreateSeaRouteAPI,
	)

	apiGroup.GET(
		"/sea-routes/draft",
		hand.GetDraftAPI,
	)

	apiGroup.PUT(
		"/sea-routes/draft",
		hand.PublishSeaRouteAPI,
	)

	apiGroup.GET(
		"/sea-routes/feed",
		hand.GetSeaRouteFeedAPI,
	)

	apiGroup.POST(
		"/sea-routes/:id/likes",
		hand.SetSeaRouteLikeAPI,
	)

	apiGroup.DELETE(
		"/sea-routes/:id",
		hand.DeleteSeaRouteAPI,
	)

	apiGroup.POST(
		"/users",
		hand.RegisterUserAPI,
	)

	apiGroup.POST(
		"/auth",
		hand.AuthAPI,
	)

	apiGroup.POST(
		"/logout",
		hand.LogoutAPI,
	)

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
