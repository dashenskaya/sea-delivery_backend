package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"sea_routes/internal/app/ds"
)

const (
	SeaRouteStatusDraft     = "черновик"
	SeaRouteStatusPublished = "опубликован"
	SeaRouteStatusDeleted   = "удален"
)

type RepositorySettings struct {
	PostgresDSN string

	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
}

type Repository struct {
	db *gorm.DB

	minio       *minio.Client
	minioBucket string
}

func New(
	settings RepositorySettings,
) (*Repository, error) {

	db, err :=
		gorm.Open(
			postgres.Open(settings.PostgresDSN),
			&gorm.Config{},
		)

	if err != nil {
		return nil, fmt.Errorf(
			"ошибка подключения к PostgreSQL: %w",
			err,
		)
	}

	minioClient, err :=
		minio.New(
			settings.MinioEndpoint,
			&minio.Options{
				Creds: credentials.NewStaticV4(
					settings.MinioAccessKey,
					settings.MinioSecretKey,
					"",
				),

				Secure: settings.MinioUseSSL,
			},
		)

	if err != nil {
		return nil, fmt.Errorf(
			"ошибка создания клиента MinIO: %w",
			err,
		)
	}

	exists, err :=
		minioClient.BucketExists(
			context.Background(),
			settings.MinioBucket,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"ошибка подключения к MinIO: %w",
			err,
		)
	}

	if !exists {
		return nil, fmt.Errorf(
			"бакет MinIO %q не найден",
			settings.MinioBucket,
		)
	}

	return &Repository{
		db:          db,
		minio:       minioClient,
		minioBucket: settings.MinioBucket,
	}, nil
}

// Подсчет лайков одного маршрута.
func (r *Repository) addLikesCount(
	seaRoute *ds.SeaRoute,
) error {
	var count int64

	err :=
		r.db.
			Model(&ds.SeaRouteLike{}).
			Where(
				"sea_route_id = ?",
				seaRoute.ID,
			).
			Count(&count).
			Error

	if err != nil {
		return err
	}

	seaRoute.LikesCount = count

	return nil
}

// Подсчет лайков списка маршрутов.
func (r *Repository) addLikesCounts(
	seaRoutes []ds.SeaRoute,
) error {
	for i := range seaRoutes {
		if err :=
			r.addLikesCount(
				&seaRoutes[i],
			); err != nil {

			return err
		}
	}

	return nil
}

// Первый опубликованный маршрут.
func (r *Repository) GetFirstPublishedSeaRoute() (
	*ds.SeaRoute,
	error,
) {
	var seaRoute ds.SeaRoute

	err :=
		r.db.
			Where(
				"sea_route_status = ?",
				SeaRouteStatusPublished,
			).
			Order(
				"sea_route_id ASC",
			).
			First(&seaRoute).
			Error

	if errors.Is(
		err,
		gorm.ErrRecordNotFound,
	) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if err :=
		r.addLikesCount(
			&seaRoute,
		); err != nil {

		return nil, err
	}

	return &seaRoute, nil
}

// Опубликованный маршрут по ID.
func (r *Repository) GetPublishedSeaRouteByID(
	id uint,
) (*ds.SeaRoute, error) {

	var seaRoute ds.SeaRoute

	err :=
		r.db.
			Where(
				"sea_route_id = ? AND sea_route_status = ?",
				id,
				SeaRouteStatusPublished,
			).
			First(&seaRoute).
			Error

	if errors.Is(
		err,
		gorm.ErrRecordNotFound,
	) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if err :=
		r.addLikesCount(
			&seaRoute,
		); err != nil {

		return nil, err
	}

	return &seaRoute, nil
}

// Следующий опубликованный маршрут.
func (r *Repository) GetNextPublishedSeaRoute(
	currentID uint,
) (*ds.SeaRoute, error) {

	var seaRoute ds.SeaRoute

	err :=
		r.db.
			Where(
				"sea_route_status = ? AND sea_route_id > ?",
				SeaRouteStatusPublished,
				currentID,
			).
			Order(
				"sea_route_id ASC",
			).
			First(&seaRoute).
			Error

	if errors.Is(
		err,
		gorm.ErrRecordNotFound,
	) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if err :=
		r.addLikesCount(
			&seaRoute,
		); err != nil {

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

	err :=
		r.db.
			Where(
				`sea_route_status = ?
				AND sea_route_distance_nautical_miles >= ?
				AND sea_route_distance_nautical_miles <= ?`,
				SeaRouteStatusPublished,
				minDistance,
				maxDistance,
			).
			Order(
				"sea_route_id ASC",
			).
			Find(&seaRoutes).
			Error

	if err != nil {
		return nil, err
	}

	if err :=
		r.addLikesCounts(
			seaRoutes,
		); err != nil {

		return nil, err
	}

	return seaRoutes, nil
}

// Получить черновик пользователя.
func (r *Repository) GetDraftByUser(
	userID uint,
) (*ds.SeaRoute, error) {

	var seaRoute ds.SeaRoute

	err :=
		r.db.
			Where(
				"sea_route_creator_id = ? AND sea_route_status = ?",
				userID,
				SeaRouteStatusDraft,
			).
			First(&seaRoute).
			Error

	if errors.Is(
		err,
		gorm.ErrRecordNotFound,
	) {
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

	existingDraft, err :=
		r.GetDraftByUser(userID)

	if err != nil {
		return nil, err
	}

	// У одного пользователя может быть
	// только один черновик.
	if existingDraft != nil {
		return existingDraft, nil
	}

	seaRoute := ds.SeaRoute{
		Name: name,

		// Эти поля пользователь заполнит
		// только после кнопки "Далее".
		DistanceNauticalMiles: nil,
		AveragePortDelayHours: nil,
		Description:           nil,
		FormedAt:              nil,

		ImageURL: "/static/media/default-route.jpg",
		VideoURL: "/static/media/default-route.mp4",

		Status:    SeaRouteStatusDraft,
		CreatorID: userID,
	}

	err =
		r.db.
			Create(&seaRoute).
			Error

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

	draft, err :=
		r.GetDraftByUser(userID)

	if err != nil {
		return nil, err
	}

	if draft == nil {
		return nil,
			fmt.Errorf(
				"черновик морского маршрута не найден",
			)
	}

	formedAt := time.Now()

	err =
		r.db.
			Model(&ds.SeaRoute{}).
			Where(
				"sea_route_id = ? AND sea_route_creator_id = ? AND sea_route_status = ?",
				draft.ID,
				userID,
				SeaRouteStatusDraft,
			).
			Updates(
				map[string]interface{}{
					"sea_route_distance_nautical_miles": distanceNauticalMiles,

					"sea_route_average_port_delay_hours": averagePortDelayHours,

					"sea_route_description": description,

					"sea_route_status": SeaRouteStatusPublished,

					"sea_route_formed_at": formedAt,
				},
			).
			Error

	if err != nil {
		return nil, err
	}

	return r.GetPublishedSeaRouteByID(
		draft.ID,
	)
}

// Логическое удаление маршрута текущего пользователя.
func (r *Repository) DeleteSeaRoute(
	ctx context.Context,
	seaRouteID uint,
	userID uint,
) error {

	sqlDB, err :=
		r.db.DB()

	if err != nil {
		return err
	}

	result, err :=
		sqlDB.ExecContext(
			ctx,
			`
			UPDATE sea_routes
			SET sea_route_status = $1
			WHERE sea_route_id = $2
			  AND sea_route_creator_id = $3
			  AND sea_route_status <> $1
			`,
			SeaRouteStatusDeleted,
			seaRouteID,
			userID,
		)

	if err != nil {
		return err
	}

	rowsAffected, err :=
		result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf(
			"морской маршрут не найден или не принадлежит текущему пользователю",
		)
	}

	return nil
}

// Сохранение имен изображения и видео в черновике.
func (r *Repository) UpdateDraftMedia(
	seaRouteID uint,
	userID uint,
	imageFileName string,
	videoFileName string,
) error {

	result :=
		r.db.
			Model(&ds.SeaRoute{}).
			Where(
				`sea_route_id = ?
				AND sea_route_creator_id = ?
				AND sea_route_status = ?`,
				seaRouteID,
				userID,
				SeaRouteStatusDraft,
			).
			Updates(
				map[string]interface{}{
					"sea_route_image_url": imageFileName,
					"sea_route_video_url": videoFileName,
				},
			)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf(
			"черновик морского маршрута не найден",
		)
	}

	return nil
}

// Установка или отмена лайка текущего пользователя.
// value = 1 — поставить лайк.
// value = 0 — убрать лайк.
func (r *Repository) SetSeaRouteLike(
	userID uint,
	seaRouteID uint,
	value int,
) (*ds.SeaRoute, error) {

	// Лайкать можно только опубликованный маршрут.
	seaRoute, err :=
		r.GetPublishedSeaRouteByID(
			seaRouteID,
		)

	if err != nil {
		return nil, err
	}

	if seaRoute == nil {
		return nil,
			fmt.Errorf(
				"опубликованный морской маршрут не найден",
			)
	}

	switch value {
	case 1:
		// FirstOrCreate не позволит повторно
		// создать тот же лайк этого пользователя.
		like := ds.SeaRouteLike{
			UserID:     userID,
			SeaRouteID: seaRouteID,
		}

		err =
			r.db.
				Where(
					"sea_route_user_id = ? AND sea_route_id = ?",
					userID,
					seaRouteID,
				).
				FirstOrCreate(&like).
				Error

		if err != nil {
			return nil, err
		}

	case 0:
		err =
			r.db.
				Where(
					"sea_route_user_id = ? AND sea_route_id = ?",
					userID,
					seaRouteID,
				).
				Delete(
					&ds.SeaRouteLike{},
				).
				Error

		if err != nil {
			return nil, err
		}

	default:
		return nil,
			fmt.Errorf(
				"значение лайка должно быть 0 или 1",
			)
	}

	// Получаем маршрут заново,
	// чтобы LikesCount пересчитался.
	return r.GetPublishedSeaRouteByID(
		seaRouteID,
	)
}

// Проверка, поставил ли текущий пользователь лайк маршруту.
func (r *Repository) IsSeaRouteLikedByUser(
	userID uint,
	seaRouteID uint,
) (bool, error) {

	var count int64

	err :=
		r.db.
			Model(&ds.SeaRouteLike{}).
			Where(
				"sea_route_user_id = ? AND sea_route_id = ?",
				userID,
				seaRouteID,
			).
			Count(&count).
			Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
