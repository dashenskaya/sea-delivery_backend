package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"sea_routes/internal/app/ds"
)

const (
	SeaRouteStatusDraft     = "черновик"
	SeaRouteStatusPublished = "опубликован"
	SeaRouteStatusDeleted   = "удален"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

// Подсчет лайков одного маршрута.
func (r *Repository) addLikesCount(seaRoute *ds.SeaRoute) error {
	var count int64

	err := r.db.
		Model(&ds.SeaRouteLike{}).
		Where("sea_route_id = ?", seaRoute.ID).
		Count(&count).
		Error

	if err != nil {
		return err
	}

	seaRoute.LikesCount = count

	return nil
}

// Подсчет лайков списка маршрутов.
func (r *Repository) addLikesCounts(seaRoutes []ds.SeaRoute) error {
	for i := range seaRoutes {
		if err := r.addLikesCount(&seaRoutes[i]); err != nil {
			return err
		}
	}

	return nil
}

// Первый опубликованный маршрут.
func (r *Repository) GetFirstPublishedSeaRoute() (*ds.SeaRoute, error) {
	var seaRoute ds.SeaRoute

	err := r.db.
		Where("sea_route_status = ?", SeaRouteStatusPublished).
		Order("sea_route_id ASC").
		First(&seaRoute).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if err := r.addLikesCount(&seaRoute); err != nil {
		return nil, err
	}

	return &seaRoute, nil
}

// Опубликованный маршрут по ID.
func (r *Repository) GetPublishedSeaRouteByID(
	id uint,
) (*ds.SeaRoute, error) {

	var seaRoute ds.SeaRoute

	err := r.db.
		Where(
			"sea_route_id = ? AND sea_route_status = ?",
			id,
			SeaRouteStatusPublished,
		).
		First(&seaRoute).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if err := r.addLikesCount(&seaRoute); err != nil {
		return nil, err
	}

	return &seaRoute, nil
}

// Следующий опубликованный маршрут.
func (r *Repository) GetNextPublishedSeaRoute(
	currentID uint,
) (*ds.SeaRoute, error) {

	var seaRoute ds.SeaRoute

	err := r.db.
		Where(
			"sea_route_status = ? AND sea_route_id > ?",
			SeaRouteStatusPublished,
			currentID,
		).
		Order("sea_route_id ASC").
		First(&seaRoute).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if err := r.addLikesCount(&seaRoute); err != nil {
		return nil, err
	}

	return &seaRoute, nil
}

// Каталог с фильтрацией по расстоянию.
func (r *Repository) GetPublishedSeaRoutesByDistance(
	minDistance int,
	maxDistance int,
) ([]ds.SeaRoute, error) {

	var seaRoutes []ds.SeaRoute

	err := r.db.
		Where(
			`sea_route_status = ?
			AND sea_route_distance_nautical_miles >= ?
			AND sea_route_distance_nautical_miles <= ?`,
			SeaRouteStatusPublished,
			minDistance,
			maxDistance,
		).
		Order("sea_route_id ASC").
		Find(&seaRoutes).
		Error

	if err != nil {
		return nil, err
	}

	if err := r.addLikesCounts(seaRoutes); err != nil {
		return nil, err
	}

	return seaRoutes, nil
}

// Получить черновик пользователя.
func (r *Repository) GetDraftByUser(
	userID uint,
) (*ds.SeaRoute, error) {

	var seaRoute ds.SeaRoute

	err := r.db.
		Where(
			"sea_route_creator_id = ? AND sea_route_status = ?",
			userID,
			SeaRouteStatusDraft,
		).
		First(&seaRoute).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &seaRoute, nil
}

// Создание черновика через ORM.
func (r *Repository) CreateDraft(
	userID uint,
	name string,
) (*ds.SeaRoute, error) {

	existingDraft, err := r.GetDraftByUser(userID)
	if err != nil {
		return nil, err
	}

	if existingDraft != nil {
		return existingDraft, nil
	}

	seaRoute := ds.SeaRoute{
		Name:                  name,
		DeparturePort:         "",
		ArrivalPort:           "",
		DistanceNauticalMiles: 0,
		AveragePortDelayHours: 0,
		Description:           "",
		ImageURL:              "",
		VideoURL:              "",
		Status:                SeaRouteStatusDraft,
		CreatorID:             userID,
	}

	err = r.db.Create(&seaRoute).Error
	if err != nil {
		return nil, err
	}

	return &seaRoute, nil
}

// Публикация черновика через ORM.
func (r *Repository) PublishDraft(
	userID uint,
	distanceNauticalMiles int,
	averagePortDelayHours int,
	description string,
) (*ds.SeaRoute, error) {

	draft, err := r.GetDraftByUser(userID)
	if err != nil {
		return nil, err
	}

	if draft == nil {
		return nil, fmt.Errorf("черновик морского маршрута не найден")
	}

	formedAt := time.Now()

	err = r.db.
		Model(&ds.SeaRoute{}).
		Where(
			"sea_route_id = ? AND sea_route_status = ?",
			draft.ID,
			SeaRouteStatusDraft,
		).
		Updates(map[string]interface{}{
			"sea_route_distance_nautical_miles":  distanceNauticalMiles,
			"sea_route_average_port_delay_hours": averagePortDelayHours,
			"sea_route_description":              description,
			"sea_route_status":                   SeaRouteStatusPublished,
			"sea_route_formed_at":                formedAt,
		}).
		Error

	if err != nil {
		return nil, err
	}

	return r.GetPublishedSeaRouteByID(draft.ID)
}

// Логическое удаление через чистый SQL UPDATE.
func (r *Repository) DeleteSeaRoute(
	ctx context.Context,
	seaRouteID uint,
) error {

	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	result, err := sqlDB.ExecContext(
		ctx,
		`
		UPDATE sea_routes
		SET sea_route_status = $1
		WHERE sea_route_id = $2
		`,
		SeaRouteStatusDeleted,
		seaRouteID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf(
			"морской маршрут с id %d не найден",
			seaRouteID,
		)
	}

	return nil
}
