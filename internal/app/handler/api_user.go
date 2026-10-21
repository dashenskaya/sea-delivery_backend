package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"sea_routes/internal/app/serializer"
)

type RegisterUserRequest struct {
	Login string `json:"login"`

	Password string `json:"password"`
}

// POST /api/users
// Регистрация нового пользователя.

// RegisterUserAPI godoc
// @Summary Регистрация пользователя
// @Description Создание нового пользователя Sea Routes
// @Tags Авторизация
// @Accept json
// @Produce json
// @Param request body RegisterUserRequest true "Логин и пароль"
// @Success 201 {object} serializer.User
// @Failure 400
// @Failure 409
// @Router /users [post]
func (h *Handler) RegisterUserAPI(
	ctx *gin.Context,
) {
	var request RegisterUserRequest

	if err :=
		ctx.ShouldBindJSON(
			&request,
		); err != nil {

		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	login :=
		strings.TrimSpace(
			request.Login,
		)

	password :=
		request.Password

	if login == "" {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	if len(login) > 25 {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	if strings.TrimSpace(password) == "" ||
		len(password) < 8 ||
		len(password) > 72 {

		ctx.Status(http.StatusBadRequest)
		return
	}

	user, err :=
		h.Repository.CreateUser(
			login,
			password,
		)

	if err != nil {
		ctx.Status(
			http.StatusConflict,
		)
		return
	}

	ctx.JSON(
		http.StatusCreated,
		serializer.FromUser(
			user,
		),
	)
}
