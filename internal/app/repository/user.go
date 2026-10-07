package repository

import (
	"fmt"

	"sea_routes/internal/app/ds"
)

// Создание нового пользователя.
func (r *Repository) CreateUser(
	login string,
	password string,
) (*ds.SeaRouteUser, error) {

	var count int64

	err :=
		r.db.
			Model(&ds.SeaRouteUser{}).
			Where(
				"sea_route_user_login = ?",
				login,
			).
			Count(&count).
			Error

	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil,
			fmt.Errorf(
				"пользователь с таким логином уже существует",
			)
	}

	user := ds.SeaRouteUser{
		Login:    login,
		Password: password,
	}

	err =
		r.db.
			Create(&user).
			Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
