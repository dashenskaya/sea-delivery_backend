package redisstore

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const blacklistPrefix = "sea_routes.jwt_blacklist."

type Store struct {
	client *redis.Client
}

// Создание подключения к Redis.
func New(addr string) *Store {
	return &Store{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
			DB:   0,
		}),
	}
}

// Проверка доступности Redis.
func (s *Store) Ping(ctx context.Context) error {
	return s.client.Ping(ctx).Err()
}

// Закрытие соединения с Redis.
func (s *Store) Close() error {
	return s.client.Close()
}

// Формирование ключа для blacklist.
// Вместо самого JWT храним его SHA-256 хеш.
func tokenKey(token string) string {
	hash := sha256.Sum256([]byte(token))

	return blacklistPrefix +
		hex.EncodeToString(hash[:])
}

// Добавление JWT в blacklist после выхода.
func (s *Store) Revoke(
	ctx context.Context,
	token string,
	expiresAt time.Time,
) error {

	ttl := time.Until(expiresAt)

	if ttl <= 0 {
		return errors.New("JWT уже истек")
	}

	return s.client.Set(
		ctx,
		tokenKey(token),
		"revoked",
		ttl,
	).Err()
}

// Проверка, находится ли JWT в blacklist.
func (s *Store) IsRevoked(
	ctx context.Context,
	token string,
) (bool, error) {

	count, err := s.client.Exists(
		ctx,
		tokenKey(token),
	).Result()

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
