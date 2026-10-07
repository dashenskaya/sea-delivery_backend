package repository

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

// UploadMediaFile загружает файл в MinIO
// и возвращает сгенерированное имя файла.
//
// kind должен быть "image" или "video".
func (r *Repository) UploadMediaFile(
	ctx context.Context,
	header *multipart.FileHeader,
	kind string,
) (string, error) {

	file, err := header.Open()
	if err != nil {
		return "", fmt.Errorf(
			"не удалось открыть файл: %w",
			err,
		)
	}

	defer file.Close()

	// Определяем реальный тип файла
	// по первым 512 байтам.
	buffer := make([]byte, 512)

	count, err := file.Read(buffer)
	if err != nil {
		return "", fmt.Errorf(
			"не удалось прочитать файл: %w",
			err,
		)
	}

	contentType :=
		http.DetectContentType(
			buffer[:count],
		)

	// Возвращаем указатель чтения
	// в начало файла перед загрузкой.
	_, err = file.Seek(0, 0)

	if err != nil {
		return "", fmt.Errorf(
			"не удалось обработать файл: %w",
			err,
		)
	}

	// Проверяем тип.
	switch kind {
	case "image":
		if !strings.HasPrefix(
			contentType,
			"image/",
		) {
			return "",
				fmt.Errorf(
					"файл должен быть изображением",
				)
		}

	case "video":
		if !strings.HasPrefix(
			contentType,
			"video/",
		) {
			return "",
				fmt.Errorf(
					"файл должен быть видео",
				)
		}

	default:
		return "",
			fmt.Errorf(
				"неизвестный тип файла",
			)
	}

	// Берём расширение исходного файла,
	// но само имя генерируем на латинице.
	extension :=
		strings.ToLower(
			filepath.Ext(
				header.Filename,
			),
		)

	fileName :=
		fmt.Sprintf(
			"%s_%s%s",
			kind,
			uuid.NewString(),
			extension,
		)

	_, err =
		r.minio.PutObject(
			ctx,
			r.minioBucket,
			fileName,
			file,
			header.Size,
			minio.PutObjectOptions{
				ContentType: contentType,
			},
		)

	if err != nil {
		return "",
			fmt.Errorf(
				"не удалось загрузить файл в MinIO: %w",
				err,
			)
	}

	return fileName, nil
}
