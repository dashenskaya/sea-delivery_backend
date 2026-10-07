package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// POST /api/auth
// Заглушка аутентификации для ЛР4.
func (h *Handler) AuthAPI(
	ctx *gin.Context,
) {
	var request AuthRequest

	if err :=
		ctx.ShouldBindJSON(
			&request,
		); err != nil {

		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	if strings.TrimSpace(request.Login) == "" ||
		strings.TrimSpace(request.Password) == "" {

		ctx.Status(
			http.StatusBadRequest,
		)
		return
	}

	ctx.Status(
		http.StatusOK,
	)
}

// POST /api/logout
// Заглушка деавторизации для ЛР4.
func (h *Handler) LogoutAPI(
	ctx *gin.Context,
) {
	ctx.Status(
		http.StatusOK,
	)
}
