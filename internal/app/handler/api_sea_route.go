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

// GetSeaRoutesAPI godoc
// @Summary Получение списка морских маршрутов
// @Description Возвращает опубликованные маршруты с фильтрацией по расстоянию
// @Tags Морские маршруты
// @Produce json
// @Param min_distance_nautical_miles query int false "Минимальное расстояние"
// @Param max_distance_nautical_miles query int false "Максимальное расстояние"
// @Success 200 {array} serializer.SeaRoute
// @Failure 400
// @Failure 500
// @Router /sea-routes [get]
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
			ctx.GetUint("user_id"),
		)

	ctx.JSON(
		http.StatusOK,
		response,
	)
}

// POST /api/sea-routes
// Создание нового черновика с изображением и видео.

// CreateSeaRouteAPI godoc
// @Summary Создание черновика морского маршрута
// @Description Создание маршрута с фотографией и видео
// @Tags Морские маршруты
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Название маршрута"
// @Param image formData file true "Фотография маршрута"
// @Param video formData file true "Видео маршрута"
// @Success 201 {object} serializer.SeaRoute
// @Failure 400
// @Failure 401
// @Failure 409
// @Failure 500
// @Security BearerAuth
// @Router /sea-routes [post]
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
			ctx.GetUint("user_id"),
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
			ctx.GetUint("user_id"),
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
			ctx.GetUint("user_id"),
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
			ctx.GetUint("user_id"),
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
			ctx.GetUint("user_id"),
		),
	)
}

// GET /api/sea-routes/draft
// Получение черновика текущего пользователя.

// GetDraftAPI godoc
// @Summary Получение черновика морского маршрута
// @Description Получение черновика текущего пользователя
// @Tags Морские маршруты
// @Produce json
// @Success 200 {object} serializer.SeaRoute
// @Failure 401
// @Failure 404
// @Failure 500
// @Security BearerAuth
// @Router /sea-routes/draft [get]
func (h *Handler) GetDraftAPI(
	ctx *gin.Context,
) {
	draft, err :=
		h.Repository.GetDraftByUser(
			ctx.GetUint("user_id"),
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
			ctx.GetUint("user_id"),
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

// PublishSeaRouteAPI godoc
// @Summary Публикация морского маршрута
// @Description Заполнение параметров и публикация черновика текущего пользователя
// @Tags Морские маршруты
// @Accept json
// @Produce json
// @Param request body PublishSeaRouteRequest true "Параметры морского маршрута"
// @Success 200 {object} serializer.SeaRoute
// @Failure 400
// @Failure 401
// @Failure 404
// @Security BearerAuth
// @Router /sea-routes/draft [put]
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
			ctx.GetUint("user_id"),
			*request.DistanceNauticalMiles,
			*request.AveragePortDelayHours,
			description,
		)

	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(
		http.StatusOK,
		serializer.FromSeaRoute(
			seaRoute,
			ctx.GetUint("user_id"),
		),
	)
}

// GET /api/sea-routes/feed
// Получение ленты опубликованных морских маршрутов.

// GetSeaRouteFeedAPI godoc
// @Summary Получение ленты морских маршрутов
// @Description Получение первого, выбранного или следующего опубликованного маршрута
// @Tags Морские маршруты
// @Produce json
// @Param sea_route_id query int false "ID текущего маршрута"
// @Param next query bool false "Перейти к следующему маршруту"
// @Success 200 {object} serializer.FeedSeaRoute
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /sea-routes/feed [get]
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
				ctx.GetUint("user_id"),
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
			ctx.GetUint("user_id"),
			isLiked,
		),
	)
}

type SeaRouteLikeRequest struct {
	Value *int `json:"value"`
}

// POST /api/sea-routes/:id/likes
// Установка или отмена лайка текущего пользователя.

// SetSeaRouteLikeAPI godoc
// @Summary Установка или отмена лайка
// @Description Значение value=1 устанавливает лайк, value=0 отменяет
// @Tags Морские маршруты
// @Accept json
// @Produce json
// @Param id path int true "ID морского маршрута"
// @Param request body SeaRouteLikeRequest true "Значение лайка"
// @Success 200 {object} serializer.SeaRoute
// @Failure 400
// @Failure 401
// @Failure 404
// @Security BearerAuth
// @Router /sea-routes/{id}/likes [post]
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
			ctx.GetUint("user_id"),
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
			ctx.GetUint("user_id"),
		),
	)
}

// DELETE /api/sea-routes/:id
// Логическое удаление маршрута текущего пользователя.

// DeleteSeaRouteAPI godoc
// @Summary Удаление морского маршрута
// @Description Логическое удаление собственного маршрута
// @Tags Морские маршруты
// @Param id path int true "ID морского маршрута"
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 404
// @Security BearerAuth
// @Router /sea-routes/{id} [delete]
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
			ctx.GetUint("user_id"),
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
