package api

import (
	"log"

	"github.com/gin-gonic/gin"

	"sea_routes/internal/app/handler"
	"sea_routes/internal/app/redisstore"
	"sea_routes/internal/app/repository"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "sea_routes/docs"
)

func StartServer(
	rep *repository.Repository,
	redisStore *redisstore.Store,
) {
	log.Println("Server start up")

	router := gin.Default()

	router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	router.LoadHTMLGlob("templates/*")

	router.Static(
		"/static",
		"./resources",
	)

	hand := handler.NewHandler(rep)

	apiGroup := router.Group("/api")

	// Публичный список морских маршрутов.
	apiGroup.GET(
		"/sea-routes",
		hand.GetSeaRoutesAPI,
	)

	// Создание маршрута — только с JWT.
	apiGroup.POST(
		"/sea-routes",
		handler.RequireAuth(redisStore),
		hand.CreateSeaRouteAPI,
	)

	// Получение черновика — только с JWT.
	apiGroup.GET(
		"/sea-routes/draft",
		handler.RequireAuth(redisStore),
		hand.GetDraftAPI,
	)

	// Публикация маршрута — только с JWT.
	apiGroup.PUT(
		"/sea-routes/draft",
		handler.RequireAuth(redisStore),
		hand.PublishSeaRouteAPI,
	)

	// Публичная лента маршрутов.
	apiGroup.GET(
		"/sea-routes/feed",
		hand.GetSeaRouteFeedAPI,
	)

	// Лайк — только с JWT.
	apiGroup.POST(
		"/sea-routes/:id/likes",
		handler.RequireAuth(redisStore),
		hand.SetSeaRouteLikeAPI,
	)

	// Удаление маршрута — только с JWT.
	apiGroup.DELETE(
		"/sea-routes/:id",
		handler.RequireAuth(redisStore),
		hand.DeleteSeaRouteAPI,
	)

	// Публичная регистрация.
	apiGroup.POST(
		"/users",
		hand.RegisterUserAPI,
	)

	// Публичная аутентификация.
	apiGroup.POST(
		"/auth",
		hand.AuthAPI,
	)

	// Выход пользователя.
	apiGroup.POST(
		"/logout",
		handler.RequireAuth(redisStore),
		hand.LogoutAPI,
	)

	// HTML-страницы из предыдущих лабораторных.
	router.GET(
		"/sea_routes/feed",
		hand.SeaRouteFeedHandler,
	)

	router.GET(
		"/sea_routes/add",
		handler.RequireAuth(redisStore),
		hand.SeaRouteAddHandler,
	)

	router.GET(
		"/sea_routes/catalog",
		hand.SeaRouteCatalogHandler,
	)

	router.POST(
		"/sea_routes/create",
		handler.RequireAuth(redisStore),
		hand.SeaRouteCreateHandler,
	)

	router.POST(
		"/sea_routes/publish",
		handler.RequireAuth(redisStore),
		hand.SeaRoutePublishHandler,
	)

	router.POST(
		"/sea_routes/delete",
		handler.RequireAuth(redisStore),
		hand.SeaRouteDeleteHandler,
	)

	err := router.Run()

	if err != nil {
		log.Println(err)
	}

	log.Println("Server down")
}
