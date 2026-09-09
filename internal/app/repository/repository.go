package repository

const (
	SeaRouteStatusDraft     = "черновик"
	SeaRouteStatusPublished = "опубликован"
	SeaRouteStatusDeleted   = "удален"
)

type SeaRoute struct {
	ID int

	DeparturePort string
	ArrivalPort   string

	DistanceNauticalMiles int
	AveragePortDelayHours int

	Description string

	ImageURL string
	VideoURL string

	Status string

	LikedUserIDs []int
	LikesCount   int
}

var SeaRoutes = []SeaRoute{
	{
		ID:                    1,
		DeparturePort:         "Шанхай",
		ArrivalPort:           "Роттердам",
		DistanceNauticalMiles: 10664,
		AveragePortDelayHours: 156,
		Description:           "Морской контейнерный маршрут из Китая в Нидерланды через крупнейшие торговые пути между Азией и Европой.",
		ImageURL:              "http://localhost:9000/sea-routes/shanghai-rotterdam.jpg",
		VideoURL:              "http://localhost:9000/sea-routes/shanghai-rotterdam.mp4",
		Status:                SeaRouteStatusPublished,
		LikedUserIDs:          []int{1, 3, 5, 8, 12},
	},
	{
		ID:                    2,
		DeparturePort:         "Гданьск",
		ArrivalPort:           "Нью-Йорк",
		DistanceNauticalMiles: 4315,
		AveragePortDelayHours: 233,
		Description:           "Морской контейнерный маршрут из Польши в США через Атлантический океан.",
		ImageURL:              "http://localhost:9000/sea-routes/gdansk-new-york.jpg",
		VideoURL:              "http://localhost:9000/sea-routes/gdansk-new-york.mp4",
		Status:                SeaRouteStatusPublished,
		LikedUserIDs:          []int{2, 7},
	},
	{
		ID:                    3,
		DeparturePort:         "Сингапур",
		ArrivalPort:           "Гамбург",
		DistanceNauticalMiles: 12298,
		AveragePortDelayHours: 163,
		Description:           "Морской контейнерный маршрут из Сингапура в Германию, соединяющий Юго-Восточную Азию с Северной Европой.",
		ImageURL:              "http://localhost:9000/sea-routes/singapore-hamburg.jpg",
		VideoURL:              "http://localhost:9000/sea-routes/singapore-hamburg.mp4",
		Status:                SeaRouteStatusPublished,
		LikedUserIDs:          []int{1, 2, 3, 5, 7, 9, 11, 15},
	},
	{
		ID:                    4,
		DeparturePort:         "Дубай",
		ArrivalPort:           "Мумбаи",
		DistanceNauticalMiles: 1329,
		AveragePortDelayHours: 237,
		Description:           "Морской контейнерный маршрут между портом Джебель-Али в Дубае и портом Нава-Шева в районе Мумбаи.",
		ImageURL:              "http://localhost:9000/sea-routes/dubai-mumbai.jpg",
		VideoURL:              "http://localhost:9000/sea-routes/dubai-mumbai.mp4",
		Status:                SeaRouteStatusPublished,
		LikedUserIDs:          []int{4, 6, 8, 10, 12},
	},
	{
		ID:                    5,
		DeparturePort:         "Пусан",
		ArrivalPort:           "Ванкувер",
		DistanceNauticalMiles: 4672,
		AveragePortDelayHours: 182,
		Description:           "Транстихоокеанский контейнерный маршрут из Южной Кореи в Канаду.",
		ImageURL:              "http://localhost:9000/sea-routes/busan-vancouver.jpg",
		VideoURL:              "http://localhost:9000/sea-routes/busan-vancouver.mp4",
		Status:                SeaRouteStatusPublished,
		LikedUserIDs:          []int{2, 4, 6, 9},
	},
	{
		ID:                    6,
		DeparturePort:         "Сантус",
		ArrivalPort:           "Кейптаун",
		DistanceNauticalMiles: 15106,
		AveragePortDelayHours: 144,
		Description:           "Морской контейнерный маршрут между Бразилией и Южной Африкой через Атлантический океан.",
		ImageURL:              "http://localhost:9000/sea-routes/santos-capetown.jpg",
		VideoURL:              "http://localhost:9000/sea-routes/santos-capetown.mp4",
		Status:                SeaRouteStatusPublished,
		LikedUserIDs:          []int{1, 5, 7},
	},
	{
		ID:                    7,
		DeparturePort:         "Токио",
		ArrivalPort:           "Лос-Анджелес",
		DistanceNauticalMiles: 4943,
		AveragePortDelayHours: 208,
		Description:           "Прямой транстихоокеанский контейнерный маршрут из Японии на западное побережье США.",
		ImageURL:              "http://localhost:9000/sea-routes/tokyo-los-angeles.jpg",
		VideoURL:              "http://localhost:9000/sea-routes/tokyo-los-angeles.mp4",
		Status:                SeaRouteStatusPublished,
		LikedUserIDs:          []int{3, 5, 8, 10, 13, 16},
	},
	{
		ID:                    8,
		DeparturePort:         "Антверпен",
		ArrivalPort:           "Саванна",
		DistanceNauticalMiles: 4467,
		AveragePortDelayHours: 216,
		Description:           "Морской контейнерный маршрут из Антверпена в Саванну через Атлантический океан.",
		ImageURL:              "http://localhost:9000/sea-routes/antwerp-savannah.jpg",
		VideoURL:              "http://localhost:9000/sea-routes/antwerp-savannah.mp4",
		Status:                SeaRouteStatusDraft,
		LikedUserIDs:          []int{},
	},
	{
		ID:                    9,
		DeparturePort:         "Барселона",
		ArrivalPort:           "Александрия",
		DistanceNauticalMiles: 2406,
		AveragePortDelayHours: 407,
		Description:           "Средиземноморский контейнерный маршрут из Испании в Египет.",
		ImageURL:              "http://localhost:9000/sea-routes/barcelona-alexandria.jpg",
		VideoURL:              "http://localhost:9000/sea-routes/barcelona-alexandria.mp4",
		Status:                SeaRouteStatusDeleted,
		LikedUserIDs:          []int{3, 6},
	},
}
