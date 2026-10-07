package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"sea_routes/internal/app/ds"
	"sea_routes/internal/app/serializer"
)

// GET /api/sea-routes
// Список только опубликованных маршрутов
// с фильтрацией по расстоянию.
func (h *Handler) GetSeaRoutesAPI(
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

		if err != nil || value < 0 {
			ctx.Status(
				http.StatusBadRequest,
			)
			return
		}

		minDistance = value
	}

	if maxString != "" {
		value, err :=
			strconv.Atoi(maxString)

		if err != nil || value < 0 {
			ctx.Status(
				http.StatusBadRequest,
			)
			return
		}

		maxDistance = value
	}

	if minDistance > maxDistance {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	seaRoutes, err :=
		h.Repository.
			GetPublishedSeaRoutesByDistance(
				minDistance,
				maxDistance,
			)

	if err != nil {
		ctx.Status(
			http.StatusInternalServerError,
		)
		return
	}

	response :=
		serializer.FromSeaRoutes(
			seaRoutes,
			CurrentUser().ID,
		)

	ctx.JSON(
		http.StatusOK,
		response,
	)
}

// POST /api/sea-routes
// Создание нового черновика с изображением и видео.
func (h *Handler) CreateSeaRouteAPI(
	ctx *gin.Context,
) {
	err :=
		ctx.Request.ParseMultipartForm(
			32 << 20,
		)

	if err != nil {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	name :=
		strings.TrimSpace(
			ctx.PostForm("name"),
		)

	if name == "" {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	// У текущего пользователя
	// может быть только один черновик.
	existingDraft, err :=
		h.Repository.GetDraftByUser(
			CurrentUser().ID,
		)

	if err != nil {
		ctx.Status(
			http.StatusInternalServerError,
		)
		return
	}

	if existingDraft != nil {
		ctx.Status(
			http.StatusConflict,
		)
		return
	}

	// Изображение передаётся как файл.
	imageHeader, err :=
		ctx.FormFile("image")

	if err != nil {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	// Видео передаётся как файл.
	videoHeader, err :=
		ctx.FormFile("video")

	if err != nil {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	imageFileName, err :=
		h.Repository.UploadMediaFile(
			ctx.Request.Context(),
			imageHeader,
			"image",
		)

	if err != nil {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	videoFileName, err :=
		h.Repository.UploadMediaFile(
			ctx.Request.Context(),
			videoHeader,
			"video",
		)

	if err != nil {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	draft, err :=
		h.Repository.CreateDraft(
			CurrentUser().ID,
			name,
		)

	if err != nil {
		ctx.Status(
			http.StatusInternalServerError,
		)
		return
	}

	err =
		h.Repository.UpdateDraftMedia(
			draft.ID,
			CurrentUser().ID,
			imageFileName,
			videoFileName,
		)

	if err != nil {
		ctx.Status(
			http.StatusInternalServerError,
		)
		return
	}

	draft, err =
		h.Repository.GetDraftByUser(
			CurrentUser().ID,
		)

	if err != nil {
		ctx.Status(
			http.StatusInternalServerError,
		)
		return
	}

	if draft == nil {
		ctx.Status(
			http.StatusInternalServerError,
		)
		return
	}

	ctx.JSON(
		http.StatusCreated,
		serializer.FromSeaRoute(
			draft,
			CurrentUser().ID,
		),
	)
}

// GET /api/sea-routes/draft
// Получение черновика текущего пользователя.
func (h *Handler) GetDraftAPI(
	ctx *gin.Context,
) {
	draft, err :=
		h.Repository.GetDraftByUser(
			CurrentUser().ID,
		)

	if err != nil {
		ctx.Status(
			http.StatusInternalServerError,
		)
		return
	}

	if draft == nil {
		ctx.Status(
			http.StatusNotFound,
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		serializer.FromSeaRoute(
			draft,
			CurrentUser().ID,
		),
	)
}

type PublishSeaRouteRequest struct {
	DistanceNauticalMiles *int `json:"distance_nautical_miles"`

	AveragePortDelayHours *int `json:"average_port_delay_hours"`

	Description string `json:"description"`
}

// PUT /api/sea-routes/draft
// Публикация черновика текущего пользователя.
func (h *Handler) PublishSeaRouteAPI(
	ctx *gin.Context,
) {
	var request PublishSeaRouteRequest

	if err :=
		ctx.ShouldBindJSON(
			&request,
		); err != nil {

		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	if request.DistanceNauticalMiles == nil ||
		*request.DistanceNauticalMiles <= 0 {

		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	if request.AveragePortDelayHours == nil ||
		*request.AveragePortDelayHours < 0 {

		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	description :=
		strings.TrimSpace(
			request.Description,
		)

	if description == "" {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	seaRoute, err :=
		h.Repository.PublishDraft(
			CurrentUser().ID,
			*request.DistanceNauticalMiles,
			*request.AveragePortDelayHours,
			description,
		)

	if err != nil {
		ctx.Status(
			http.StatusNotFound,
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		serializer.FromSeaRoute(
			seaRoute,
			CurrentUser().ID,
		),
	)
}

// GET /api/sea-routes/feed
// Получение ленты опубликованных морских маршрутов.
func (h *Handler) GetSeaRouteFeedAPI(
	ctx *gin.Context,
) {
	seaRouteIDString :=
		ctx.Query(
			"sea_route_id",
		)

	nextString :=
		ctx.Query(
			"next",
		)

	var seaRoute *ds.SeaRoute
	var err error

	// Без ID возвращается первый
	// опубликованный маршрут.
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

		if parseErr != nil ||
			seaRouteID == 0 {

			ctx.Status(
				http.StatusBadRequest,
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
		ctx.Status(
			http.StatusInternalServerError,
		)
		return
	}

	if seaRoute == nil {
		ctx.Status(
			http.StatusNotFound,
		)
		return
	}

	// next=true означает переход
	// к следующему опубликованному маршруту.
	if nextString == "true" {
		nextSeaRoute, nextErr :=
			h.Repository.
				GetNextPublishedSeaRoute(
					seaRoute.ID,
				)

		if nextErr != nil {
			ctx.Status(
				http.StatusInternalServerError,
			)
			return
		}

		// После последнего маршрута
		// возвращаемся к первому.
		if nextSeaRoute == nil {
			nextSeaRoute, nextErr =
				h.Repository.
					GetFirstPublishedSeaRoute()

			if nextErr != nil {
				ctx.Status(
					http.StatusInternalServerError,
				)
				return
			}

			if nextSeaRoute == nil {
				ctx.Status(
					http.StatusNotFound,
				)
				return
			}
		}

		seaRoute = nextSeaRoute
	}

	// Признак 0/1:
	// лайкнул ли текущий пользователь маршрут.
	isLiked, err :=
		h.Repository.
			IsSeaRouteLikedByUser(
				CurrentUser().ID,
				seaRoute.ID,
			)

	if err != nil {
		ctx.Status(
			http.StatusInternalServerError,
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		serializer.FromSeaRouteFeed(
			seaRoute,
			CurrentUser().ID,
			isLiked,
		),
	)
}

type SeaRouteLikeRequest struct {
	Value *int `json:"value"`
}

// POST /api/sea-routes/:id/likes
// Установка или отмена лайка текущего пользователя.
func (h *Handler) SetSeaRouteLikeAPI(
	ctx *gin.Context,
) {
	idString :=
		ctx.Param(
			"id",
		)

	id, err :=
		strconv.ParseUint(
			idString,
			10,
			64,
		)

	if err != nil ||
		id == 0 {

		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	var request SeaRouteLikeRequest

	if err :=
		ctx.ShouldBindJSON(
			&request,
		); err != nil {

		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	if request.Value == nil ||
		(*request.Value != 0 &&
			*request.Value != 1) {

		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	seaRoute, err :=
		h.Repository.SetSeaRouteLike(
			CurrentUser().ID,
			uint(id),
			*request.Value,
		)

	if err != nil {
		ctx.Status(
			http.StatusNotFound,
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		serializer.FromSeaRoute(
			seaRoute,
			CurrentUser().ID,
		),
	)
}

// DELETE /api/sea-routes/:id
// Логическое удаление маршрута текущего пользователя.
func (h *Handler) DeleteSeaRouteAPI(
	ctx *gin.Context,
) {
	idString :=
		ctx.Param(
			"id",
		)

	id, err :=
		strconv.ParseUint(
			idString,
			10,
			64,
		)

	if err != nil ||
		id == 0 {

		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	err =
		h.Repository.DeleteSeaRoute(
			ctx.Request.Context(),
			uint(id),
			CurrentUser().ID,
		)

	if err != nil {
		ctx.Status(
			http.StatusNotFound,
		)
		return
	}

	ctx.Status(
		http.StatusOK,
	)
}
