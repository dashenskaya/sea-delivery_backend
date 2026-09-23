package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"sea_routes/internal/app/ds"
	"sea_routes/internal/app/repository"
)

const currentUserID uint = 1

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(
	r *repository.Repository,
) *Handler {
	return &Handler{
		Repository: r,
	}
}

// GET /sea_routes/feed
func (h *Handler) SeaRouteFeedHandler(
	ctx *gin.Context,
) {
	seaRouteIDString :=
		ctx.Query("sea_route_id")

	var seaRoute *ds.SeaRoute
	var err error

	if seaRouteIDString == "" {
		seaRoute, err =
			h.Repository.
				GetFirstPublishedSeaRoute()
	} else {
		seaRouteID, parseErr :=
			strconv.ParseUint(
				seaRouteIDString,
				10,
				64,
			)

		if parseErr != nil {
			ctx.String(
				http.StatusBadRequest,
				"Некорректный идентификатор маршрута",
			)
			return
		}

		seaRoute, err =
			h.Repository.
				GetPublishedSeaRouteByID(
					uint(seaRouteID),
				)
	}

	if err != nil {
		ctx.String(
			http.StatusInternalServerError,
			"Ошибка получения маршрута",
		)
		return
	}

	if seaRoute == nil {
		ctx.String(
			http.StatusNotFound,
			"Морской маршрут не найден",
		)
		return
	}

	// Кнопка "Следующий".
	if ctx.Query("next") == "true" {
		nextSeaRoute, nextErr :=
			h.Repository.
				GetNextPublishedSeaRoute(
					seaRoute.ID,
				)

		if nextErr != nil {
			ctx.String(
				http.StatusInternalServerError,
				"Ошибка получения следующего маршрута",
			)
			return
		}

		if nextSeaRoute != nil {
			seaRoute = nextSeaRoute
		} else {
			seaRoute, err =
				h.Repository.
					GetFirstPublishedSeaRoute()

			if err != nil {
				ctx.String(
					http.StatusInternalServerError,
					"Ошибка получения маршрута",
				)
				return
			}
		}
	}

	// Если ссылка пустая или файл недоступен,
	// здесь подставится локальный default.
	prepareSeaRouteMedia(seaRoute)

	ctx.HTML(
		http.StatusOK,
		"sea_routes_feed.html",
		gin.H{
			"seaRoute": seaRoute,
		},
	)
}

// GET /sea_routes/add
func (h *Handler) SeaRouteAddHandler(
	ctx *gin.Context,
) {
	draft, err :=
		h.Repository.
			GetDraftByUser(currentUserID)

	if err != nil {
		ctx.String(
			http.StatusInternalServerError,
			"Ошибка получения черновика: %s",
			err.Error(),
		)
		return
	}

	hasDraft := draft != nil

	seaRoute := ds.SeaRoute{}

	if draft != nil {
		seaRoute = *draft
	}

	ctx.HTML(
		http.StatusOK,
		"sea_routes_add.html",
		gin.H{
			"seaRoute": seaRoute,
			"hasDraft": hasDraft,
		},
	)
}

// GET /sea_routes/catalog
func (h *Handler) SeaRouteCatalogHandler(
	ctx *gin.Context,
) {
	minDistance := 0
	maxDistance := 16000

	minString :=
		ctx.Query(
			"min_distance_nautical_miles",
		)

	maxString :=
		ctx.Query(
			"max_distance_nautical_miles",
		)

	if minString != "" {
		value, err :=
			strconv.Atoi(minString)

		if err != nil {
			ctx.String(
				http.StatusBadRequest,
				"Некорректное минимальное расстояние",
			)
			return
		}

		minDistance = value
	}

	if maxString != "" {
		value, err :=
			strconv.Atoi(maxString)

		if err != nil {
			ctx.String(
				http.StatusBadRequest,
				"Некорректное максимальное расстояние",
			)
			return
		}

		maxDistance = value
	}

	if minDistance > maxDistance {
		minDistance, maxDistance =
			maxDistance, minDistance
	}

	seaRoutes, err :=
		h.Repository.
			GetPublishedSeaRoutesByDistance(
				minDistance,
				maxDistance,
			)

	if err != nil {
		ctx.String(
			http.StatusInternalServerError,
			"Ошибка получения каталога",
		)
		return
	}

	// Проверяем изображения перед выводом.
	prepareSeaRouteImages(seaRoutes)

	ctx.HTML(
		http.StatusOK,
		"sea_routes_catalog.html",
		gin.H{
			"seaRoutes": seaRoutes,

			"minDistanceNauticalMiles": minDistance,

			"maxDistanceNauticalMiles": maxDistance,
		},
	)
}

// POST /sea_routes/create
// Создание черновика через ORM.
func (h *Handler) SeaRouteCreateHandler(
	ctx *gin.Context,
) {
	name :=
		strings.TrimSpace(
			ctx.PostForm(
				"sea_route_name",
			),
		)

	if name == "" {
		ctx.String(
			http.StatusBadRequest,
			"Название маршрута обязательно",
		)
		return
	}

	_, err :=
		h.Repository.CreateDraft(
			currentUserID,
			name,
		)

	if err != nil {
		ctx.String(
			http.StatusInternalServerError,
			"Ошибка создания черновика: %s",
			err.Error(),
		)
		return
	}

	ctx.Redirect(
		http.StatusSeeOther,
		"/sea_routes/add",
	)
}

// POST /sea_routes/publish
// Публикация через ORM.
func (h *Handler) SeaRoutePublishHandler(
	ctx *gin.Context,
) {
	description :=
		strings.TrimSpace(
			ctx.PostForm("description"),
		)

	distance, err :=
		strconv.Atoi(
			ctx.PostForm(
				"distance_nautical_miles",
			),
		)

	if err != nil || distance < 0 {
		ctx.String(
			http.StatusBadRequest,
			"Некорректное расстояние",
		)
		return
	}

	delay, err :=
		strconv.Atoi(
			ctx.PostForm(
				"average_port_delay_hours",
			),
		)

	if err != nil || delay < 0 {
		ctx.String(
			http.StatusBadRequest,
			"Некорректная задержка",
		)
		return
	}

	if description == "" {
		ctx.String(
			http.StatusBadRequest,
			"Краткое описание обязательно",
		)
		return
	}

	seaRoute, err :=
		h.Repository.PublishDraft(
			currentUserID,
			distance,
			delay,
			description,
		)

	if err != nil {
		ctx.String(
			http.StatusInternalServerError,
			"Ошибка публикации маршрута: %s",
			err.Error(),
		)
		return
	}

	ctx.Redirect(
		http.StatusSeeOther,
		fmt.Sprintf(
			"/sea_routes/feed?sea_route_id=%d",
			seaRoute.ID,
		),
	)
}

// POST /sea_routes/delete
// Логическое удаление через SQL UPDATE.
func (h *Handler) SeaRouteDeleteHandler(
	ctx *gin.Context,
) {
	idString :=
		ctx.PostForm("sea_route_id")

	id, err :=
		strconv.ParseUint(
			idString,
			10,
			64,
		)

	if err != nil {
		ctx.String(
			http.StatusBadRequest,
			"Некорректный идентификатор маршрута",
		)
		return
	}

	err =
		h.Repository.DeleteSeaRoute(
			ctx.Request.Context(),
			uint(id),
		)

	if err != nil {
		ctx.String(
			http.StatusInternalServerError,
			"Ошибка удаления маршрута: %s",
			err.Error(),
		)
		return
	}

	ctx.Redirect(
		http.StatusSeeOther,
		"/sea_routes/catalog",
	)
}
