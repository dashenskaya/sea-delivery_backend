package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sea_delivery/internal/app/repository"
)

func SeaRouteFeedHandler(ctx *gin.Context) {
	var publishedSeaRoutes []repository.SeaRoute

	for _, seaRoute := range repository.SeaRoutes {
		if seaRoute.Status == repository.SeaRouteStatusPublished {
			seaRoute.LikesCount = len(seaRoute.LikedUserIDs)
			publishedSeaRoutes = append(publishedSeaRoutes, seaRoute)
		}
	}

	if len(publishedSeaRoutes) == 0 {
		ctx.String(http.StatusNotFound, "Опубликованные морские маршруты не найдены")
		return
	}

	currentIndex := 0
	seaRouteIDString := ctx.Query("sea_route_id")

	if seaRouteIDString != "" {
		seaRouteID, err := strconv.Atoi(seaRouteIDString)

		if err != nil {
			ctx.String(http.StatusBadRequest, "Некорректный идентификатор маршрута")
			return
		}

		found := false

		for index, seaRoute := range publishedSeaRoutes {
			if seaRoute.ID == seaRouteID {
				currentIndex = index
				found = true
				break
			}
		}

		if !found {
			ctx.String(http.StatusNotFound, "Морской маршрут не найден")
			return
		}
	}

	if ctx.Query("next") == "true" {
		currentIndex++

		if currentIndex >= len(publishedSeaRoutes) {
			currentIndex = 0
		}
	}

	ctx.HTML(
		http.StatusOK,
		"sea_route_feed.html",
		gin.H{
			"seaRoute": publishedSeaRoutes[currentIndex],
		},
	)
}

func SeaRouteAddHandler(ctx *gin.Context) {
	for _, seaRoute := range repository.SeaRoutes {
		if seaRoute.Status == repository.SeaRouteStatusDraft {
			seaRoute.LikesCount = len(seaRoute.LikedUserIDs)

			ctx.HTML(
				http.StatusOK,
				"sea_route_add.html",
				gin.H{
					"seaRoute": seaRoute,
				},
			)

			return
		}
	}

	ctx.String(http.StatusNotFound, "Черновик морского маршрута не найден")
}

func SeaRouteCatalogHandler(ctx *gin.Context) {
	var publishedSeaRoutes []repository.SeaRoute

	maxDistanceString := ctx.Query("max_distance_nautical_miles")
	maxDistanceNauticalMiles := 0

	if maxDistanceString != "" {
		value, err := strconv.Atoi(maxDistanceString)

		if err != nil {
			ctx.String(http.StatusBadRequest, "Некорректное максимальное расстояние")
			return
		}

		maxDistanceNauticalMiles = value
	}

	for _, seaRoute := range repository.SeaRoutes {
		if seaRoute.Status != repository.SeaRouteStatusPublished {
			continue
		}

		if maxDistanceString != "" &&
			seaRoute.DistanceNauticalMiles > maxDistanceNauticalMiles {
			continue
		}

		seaRoute.LikesCount = len(seaRoute.LikedUserIDs)
		publishedSeaRoutes = append(publishedSeaRoutes, seaRoute)
	}

	ctx.HTML(
		http.StatusOK,
		"sea_route_catalog.html",
		gin.H{
			"seaRoutes":                publishedSeaRoutes,
			"maxDistanceNauticalMiles": maxDistanceString,
		},
	)
}
