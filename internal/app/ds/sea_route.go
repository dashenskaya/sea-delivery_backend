package ds

import "time"

type SeaRoute struct {
	ID uint `gorm:"column:sea_route_id;primaryKey;autoIncrement"`

	Name string `gorm:"column:sea_route_name;type:varchar(180);not null;default:''"`

	DeparturePort string `gorm:"column:sea_route_departure_port;type:varchar(80);not null"`

	ArrivalPort string `gorm:"column:sea_route_arrival_port;type:varchar(80);not null"`

	DistanceNauticalMiles int `gorm:"column:sea_route_distance_nautical_miles;not null;default:0"`

	AveragePortDelayHours int `gorm:"column:sea_route_average_port_delay_hours;not null;default:0"`

	Description string `gorm:"column:sea_route_description;type:text"`

	ImageURL string `gorm:"column:sea_route_image_url;type:varchar(500)"`

	VideoURL string `gorm:"column:sea_route_video_url;type:varchar(500)"`

	Status string `gorm:"column:sea_route_status;type:varchar(20);not null;default:'черновик'"`

	CreatedAt time.Time `gorm:"column:sea_route_created_at;not null;autoCreateTime"`

	FormedAt *time.Time `gorm:"column:sea_route_formed_at"`

	CreatorID uint `gorm:"column:sea_route_creator_id;not null"`

	Creator SeaRouteUser `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`

	LikesCount int64 `gorm:"-"`
}

func (SeaRoute) TableName() string {
	return "sea_routes"
}
