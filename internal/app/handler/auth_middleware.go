package handler

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"sea_routes/internal/app/redisstore"
)

// RequireAuth проверяет JWT и Redis blacklist.
func RequireAuth(store *redisstore.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authHeader := ctx.GetHeader("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(
			strings.TrimPrefix(authHeader, "Bearer "),
		)

		if tokenString == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		secret := os.Getenv("JWT_SECRET")

		if len(secret) < 32 {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Проверяем JWT: подпись, алгоритм,
		// издателя и срок действия.
		claims := &jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			},
			jwt.WithValidMethods(
				[]string{jwt.SigningMethodHS256.Alg()},
			),
			jwt.WithIssuer("sea_routes"),
			jwt.WithExpirationRequired(),
		)

		if err != nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Получаем ID пользователя из JWT.
		userID, err := strconv.ParseUint(
			claims.Subject,
			10,
			64,
		)

		if err != nil || userID == 0 ||
			userID > uint64(^uint(0)) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Проверяем, не отозван ли JWT через Redis.
		isRevoked, err := store.IsRevoked(
			ctx.Request.Context(),
			tokenString,
		)

		if err != nil {
			ctx.AbortWithStatus(http.StatusServiceUnavailable)
			return
		}

		if isRevoked {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Сохраняем данные для обработчиков.
		ctx.Set("user_id", uint(userID))
		ctx.Set("jwt_token", tokenString)
		ctx.Set("jwt_expires_at", claims.ExpiresAt.Time)
		ctx.Set("redis_store", store)

		ctx.Next()
	}
}
