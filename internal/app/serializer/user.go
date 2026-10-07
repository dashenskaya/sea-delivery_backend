package serializer

import "sea_routes/internal/app/ds"

type User struct {
	ID uint `json:"id"`

	Login string `json:"login"`
}

func FromUser(
	user *ds.SeaRouteUser,
) User {
	return User{
		ID:    user.ID,
		Login: user.Login,
	}
}
