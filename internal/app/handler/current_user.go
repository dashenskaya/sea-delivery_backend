package handler

import "sea_routes/internal/app/ds"

const currentUserID uint = 1

var currentUserSingleton = &ds.SeaRouteUser{
	ID: currentUserID,
}

// CurrentUser возвращает единственный экземпляр
// текущего пользователя-заглушки.
func CurrentUser() *ds.SeaRouteUser {
	return currentUserSingleton
}
