package serializer

import (
	"sea_routes/internal/app/ds"
)

type SeaRoute struct {
	ID uint `json:"id"`

	Name string `json:"name"`

	DistanceNauticalMiles *int `json:"distance_nautical_miles"`

	AveragePortDelayHours *int `json:"average_port_delay_hours"`

	Description *string `json:"description"`

	Image string `json:"image"`

	Video string `json:"video"`

	LikesCount int64 `json:"likes_count"`

	IsCreator int `json:"is_creator"`
}

type FeedSeaRoute struct {
	SeaRoute

	IsLiked int `json:"is_liked"`
}

func FromSeaRoute(
	seaRoute *ds.SeaRoute,
	currentUserID uint,
) SeaRoute {

	isCreator := 0

	if seaRoute.CreatorID == currentUserID {
		isCreator = 1
	}

	return SeaRoute{
		ID: seaRoute.ID,

		Name: seaRoute.Name,

		DistanceNauticalMiles: seaRoute.DistanceNauticalMiles,

		AveragePortDelayHours: seaRoute.AveragePortDelayHours,

		Description: seaRoute.Description,

		Image: seaRoute.ImageURL,

		Video: seaRoute.VideoURL,

		LikesCount: seaRoute.LikesCount,

		IsCreator: isCreator,
	}
}

func FromSeaRouteFeed(
	seaRoute *ds.SeaRoute,
	currentUserID uint,
	isLiked bool,
) FeedSeaRoute {

	liked := 0

	if isLiked {
		liked = 1
	}

	return FeedSeaRoute{
		SeaRoute: FromSeaRoute(
			seaRoute,
			currentUserID,
		),

		IsLiked: liked,
	}
}

func FromSeaRoutes(
	seaRoutes []ds.SeaRoute,
	currentUserID uint,
) []SeaRoute {

	result :=
		make(
			[]SeaRoute,
			0,
			len(seaRoutes),
		)

	for i := range seaRoutes {
		result =
			append(
				result,
				FromSeaRoute(
					&seaRoutes[i],
					currentUserID,
				),
			)
	}

	return result
}
