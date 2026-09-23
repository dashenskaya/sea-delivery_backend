package handler

import (
	"net/http"
	"strings"
	"time"

	"sea_routes/internal/app/ds"
)

const (
	defaultImageURL = "/static/media/default-route.jpg"
	defaultVideoURL = "/static/media/default-route.mp4"
)

var mediaHTTPClient = &http.Client{
	Timeout: 1 * time.Second,
}

// Проверяет, доступен ли файл по URL.
func mediaAvailable(url string) bool {
	if strings.TrimSpace(url) == "" {
		return false
	}

	req, err := http.NewRequest(
		http.MethodHead,
		url,
		nil,
	)

	if err != nil {
		return false
	}

	resp, err := mediaHTTPClient.Do(req)
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	return resp.StatusCode >= 200 &&
		resp.StatusCode < 400
}

// Подготавливает фото и видео для страницы ленты.
func prepareSeaRouteMedia(
	seaRoute *ds.SeaRoute,
) {
	if !mediaAvailable(seaRoute.ImageURL) {
		seaRoute.ImageURL = defaultImageURL
	}

	if !mediaAvailable(seaRoute.VideoURL) {
		seaRoute.VideoURL = defaultVideoURL
	}
}

// Для каталога проверяем только картинки.
func prepareSeaRouteImages(
	seaRoutes []ds.SeaRoute,
) {
	for i := range seaRoutes {
		if !mediaAvailable(seaRoutes[i].ImageURL) {
			seaRoutes[i].ImageURL =
				defaultImageURL
		}
	}
}
