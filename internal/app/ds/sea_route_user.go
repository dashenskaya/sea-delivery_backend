package ds

type SeaRouteUser struct {
	ID uint `gorm:"column:sea_route_user_id;primaryKey;autoIncrement"`

	Name string `gorm:"column:sea_route_user_name;type:varchar(80);not null"`

	Email string `gorm:"column:sea_route_user_email;type:varchar(120);not null;unique"`
}

func (SeaRouteUser) TableName() string {
	return "sea_route_users"
}
