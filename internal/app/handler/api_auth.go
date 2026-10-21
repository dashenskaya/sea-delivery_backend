package handler

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"sea_routes/internal/app/redisstore"
)

// Данные для входа пользователя.
type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Ответ при успешной аутентификации.
type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	UserID      uint   `json:"user_id"`
	Login       string `json:"login"`
}

// POST /api/auth
// Аутентификация пользователя и выдача JWT.

// AuthAPI godoc
// @Summary Аутентификация пользователя
// @Description Проверка логина и пароля и получение JWT
// @Tags Авторизация
// @Accept json
// @Produce json
// @Param request body AuthRequest true "Логин и пароль"
// @Success 200 {object} AuthResponse
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /auth [post]
func (h *Handler) AuthAPI(ctx *gin.Context) {
	var request AuthRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	login := strings.TrimSpace(request.Login)

	if login == "" || strings.TrimSpace(request.Password) == "" {
		ctx.Status(http.StatusBadRequest)
		return
	}

	// Находим пользователя в PostgreSQL.
	user, err := h.Repository.GetUserByLogin(login)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// Проверяем пароль по сохраненному хешу bcrypt.
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(request.Password),
	)

	if err != nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	// Секретный ключ берем из переменных окружения.
	secret := os.Getenv("JWT_SECRET")

	if len(secret) < 32 {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// Создаем JWT на один час.
	now := time.Now()
	expiresIn := time.Hour

	claims := jwt.RegisteredClaims{
		Subject:  strconv.FormatUint(uint64(user.ID), 10),
		IssuedAt: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(
			now.Add(expiresIn),
		),
		Issuer: "sea_routes",
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString(
		[]byte(secret),
	)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// Отправляем токен пользователю.
	ctx.JSON(http.StatusOK, AuthResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   int64(expiresIn.Seconds()),
		UserID:      user.ID,
		Login:       user.Login,
	})
}

// POST /api/logout
// Пока заглушка: Redis подключим следующим этапом.

// LogoutAPI godoc
// @Summary Выход пользователя из системы
// @Description Отзыв JWT через blacklist Redis
// @Tags Авторизация
// @Produce json
// @Success 200
// @Failure 401
// @Failure 500
// @Security BearerAuth
// @Router /logout [post]
func (h *Handler) LogoutAPI(ctx *gin.Context) {

	// Получаем JWT, проверенный middleware.
	tokenString := ctx.GetString("jwt_token")

	// Получаем время окончания действия JWT.
	expiresValue, exists := ctx.Get("jwt_expires_at")

	if !exists {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	expiresAt, ok := expiresValue.(time.Time)

	if !ok {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// Получаем Redis из контекста.
	storeValue, exists := ctx.Get("redis_store")

	if !exists {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	store, ok := storeValue.(*redisstore.Store)

	if !ok || store == nil || tokenString == "" {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// Записываем JWT в blacklist Redis
	// до окончания срока действия токена.
	err := store.Revoke(
		ctx.Request.Context(),
		tokenString,
		expiresAt,
	)

	if err != nil {
		ctx.Status(http.StatusServiceUnavailable)
		return
	}

	// Успешный выход из аккаунта.
	ctx.Status(http.StatusOK)
}
