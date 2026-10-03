package ds

type SeaRouteUser struct {
	ID uint `gorm:"column:sea_route_user_id;primaryKey;autoIncrement"`

	Login string `gorm:"column:sea_route_user_login;type:varchar(25);not null;unique"`

	Password string `gorm:"column:sea_route_user_password;type:varchar(100);not null"`
}

func (SeaRouteUser) TableName() string {
	return "sea_route_users"
}
