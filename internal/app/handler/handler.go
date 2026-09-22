package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sea_routes/internal/app/repository"
)

func getLikesCount(likedUserIDs []int) int {
	return len(likedUserIDs)
}

func SeaRouteFeedHandler(ctx *gin.Context) {
	var firstPublishedSeaRoute repository.SeaRoute
	firstPublishedFound := false

	for _, seaRoute := range repository.SeaRoutes {
		if seaRoute.Status != repository.SeaRouteStatusPublished {
			continue
		}

		seaRoute.LikesCount = getLikesCount(seaRoute.LikedUserIDs)

		if !firstPublishedFound || seaRoute.ID < firstPublishedSeaRoute.ID {
			firstPublishedSeaRoute = seaRoute
			firstPublishedFound = true
		}
	}

	if !firstPublishedFound {
		ctx.String(http.StatusNotFound, "Опубликованные морские маршруты не найдены")
		return
	}

	seaRouteIDString := ctx.Query("sea_route_id")

	if seaRouteIDString == "" {
		ctx.HTML(
			http.StatusOK,
			"sea_routes_feed.html",
			gin.H{
				"seaRoute": firstPublishedSeaRoute,
			},
		)
		return
	}

	currentSeaRouteID, err := strconv.Atoi(seaRouteIDString)

	if err != nil {
		ctx.String(http.StatusBadRequest, "Некорректный идентификатор маршрута")
		return
	}

	var currentSeaRoute repository.SeaRoute
	currentSeaRouteFound := false

	for _, seaRoute := range repository.SeaRoutes {
		if seaRoute.Status != repository.SeaRouteStatusPublished {
			continue
		}

		if seaRoute.ID == currentSeaRouteID {
			seaRoute.LikesCount = getLikesCount(seaRoute.LikedUserIDs)
			currentSeaRoute = seaRoute
			currentSeaRouteFound = true
			break
		}
	}

	if !currentSeaRouteFound {
		ctx.String(http.StatusNotFound, "Морской маршрут не найден")
		return
	}

	if ctx.Query("next") == "true" {
		var nextSeaRoute repository.SeaRoute
		nextSeaRouteFound := false

		for _, seaRoute := range repository.SeaRoutes {
			if seaRoute.Status != repository.SeaRouteStatusPublished {
				continue
			}

			if seaRoute.ID <= currentSeaRouteID {
				continue
			}

			if !nextSeaRouteFound || seaRoute.ID < nextSeaRoute.ID {
				seaRoute.LikesCount = getLikesCount(seaRoute.LikedUserIDs)
				nextSeaRoute = seaRoute
				nextSeaRouteFound = true
			}
		}

		if nextSeaRouteFound {
			currentSeaRoute = nextSeaRoute
		} else {
			currentSeaRoute = firstPublishedSeaRoute
		}
	}

	ctx.HTML(
		http.StatusOK,
		"sea_routes_feed.html",
		gin.H{
			"seaRoute": currentSeaRoute,
		},
	)
}

func SeaRouteAddHandler(ctx *gin.Context) {
	for _, seaRoute := range repository.SeaRoutes {
		if seaRoute.Status == repository.SeaRouteStatusDraft {
			seaRoute.LikesCount = getLikesCount(seaRoute.LikedUserIDs)

			ctx.HTML(
				http.StatusOK,
				"sea_routes_add.html",
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
	minDistanceNauticalMiles := 0
	maxDistanceNauticalMiles := 16000

	minDistanceString := ctx.Query("min_distance_nautical_miles")
	maxDistanceString := ctx.Query("max_distance_nautical_miles")

	if minDistanceString != "" {
		value, err := strconv.Atoi(minDistanceString)

		if err != nil {
			ctx.String(http.StatusBadRequest, "Некорректное минимальное расстояние")
			return
		}

		minDistanceNauticalMiles = value
	}

	if maxDistanceString != "" {
		value, err := strconv.Atoi(maxDistanceString)

		if err != nil {
			ctx.String(http.StatusBadRequest, "Некорректное максимальное расстояние")
			return
		}

		maxDistanceNauticalMiles = value
	}

	if minDistanceNauticalMiles > maxDistanceNauticalMiles {
		minDistanceNauticalMiles, maxDistanceNauticalMiles =
			maxDistanceNauticalMiles, minDistanceNauticalMiles
	}

	var publishedSeaRoutes []repository.SeaRoute

	for _, seaRoute := range repository.SeaRoutes {
		if seaRoute.Status != repository.SeaRouteStatusPublished {
			continue
		}

		if seaRoute.DistanceNauticalMiles < minDistanceNauticalMiles {
			continue
		}

		if seaRoute.DistanceNauticalMiles > maxDistanceNauticalMiles {
			continue
		}

		seaRoute.LikesCount = getLikesCount(seaRoute.LikedUserIDs)

		publishedSeaRoutes = append(publishedSeaRoutes, seaRoute)
	}

	ctx.HTML(
		http.StatusOK,
		"sea_routes_catalog.html",
		gin.H{
			"seaRoutes":                publishedSeaRoutes,
			"minDistanceNauticalMiles": minDistanceNauticalMiles,
			"maxDistanceNauticalMiles": maxDistanceNauticalMiles,
		},
	)
}
