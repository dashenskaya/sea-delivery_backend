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
		strings.TrimSpace(
			request.Password,
		)

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

	if password == "" {
		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	if len(password) > 100 {
		ctx.Status(
			http.StatusBadRequest,
		)
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
