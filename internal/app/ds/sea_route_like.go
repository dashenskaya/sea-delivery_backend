package ds

type SeaRouteLike struct {
	ID uint `gorm:"column:sea_route_like_id;primaryKey;autoIncrement"`

	UserID uint `gorm:"column:sea_route_user_id;not null"`

	SeaRouteID uint `gorm:"column:sea_route_id;not null"`

	User SeaRouteUser `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`

	SeaRoute SeaRoute `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

func (SeaRouteLike) TableName() string {
	return "sea_route_likes"
}
