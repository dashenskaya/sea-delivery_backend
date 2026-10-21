package repository

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"sea_routes/internal/app/ds"
)

// Создание нового пользователя с хешированием пароля.
func (r *Repository) CreateUser(
	login string,
	password string,
) (*ds.SeaRouteUser, error) {

	var count int64

	err := r.db.
		Model(&ds.SeaRouteUser{}).
		Where("sea_route_user_login = ?", login).
		Count(&count).Error

	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, fmt.Errorf(
			"пользователь с таким логином уже существует",
		)
	}

	// Хешируем пароль перед сохранением в базу.
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user := ds.SeaRouteUser{
		Login:    login,
		Password: string(passwordHash),
	}

	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Поиск пользователя по логину для аутентификации.
func (r *Repository) GetUserByLogin(
	login string,
) (*ds.SeaRouteUser, error) {

	var user ds.SeaRouteUser

	err := r.db.
		Where("sea_route_user_login = ?", login).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
