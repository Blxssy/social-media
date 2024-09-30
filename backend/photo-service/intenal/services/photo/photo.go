package photo

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Blxssy/social-media/photo-service/intenal/models"
	"github.com/Blxssy/social-media/photo-service/intenal/storage"
	"io/ioutil"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Photo - структура, представляющая сервис работы с фотографиями.
// Поля:
// - log: используется для логирования событий.
// - storage: используется для взаимодействия с базой данных.
type Photo struct {
	log     *slog.Logger
	storage storage.Storage
}

// UploadPhoto - загружает фотографию на Яндекс Диск и сохраняет данные о фото в базе данных.
// Параметры:
// - ctx: контекст для управления запросом.
// - photoData: байтовые данные изображения.
// - filename: имя файла для сохранения на Яндекс Диске.
// - uid: ID пользователя, загружающего фото.
// Возвращает ID загруженного фото или ошибку.
func (p *Photo) UploadPhoto(ctx context.Context, photoData []byte, filename string, uid uint64) (uint64, error) {
	// Загрузка фотографии на Яндекс Диск
	yandexDiskUploadURL := "https://cloud-api.yandex.net/v1/disk/resources/upload?path=" + filename

	req, err := http.NewRequest("PUT", yandexDiskUploadURL, bytes.NewReader(photoData))
	if err != nil {
		p.log.Error("failed to create request for Yandex Disk", err)
		return 0, fmt.Errorf("failed to create request for Yandex Disk: %v", err)
	}

	authHeader := fmt.Sprintf("OAuth %s", os.Getenv("YANDEX_DISK_TOKEN"))
	fmt.Println(yandexDiskUploadURL, authHeader)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		p.log.Error("failed to upload photo to Yandex Disk", err)
		return 0, fmt.Errorf("failed to upload photo to Yandex Disk: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		p.log.Error("error uploading photo to Yandex Disk", string(body))
		return 0, fmt.Errorf("failed to upload photo: %s", body)
	}

	// Сохранение метаданных фотографии в базе данных
	imageURL := "https://disk.yandex.ru/client/disk/" + filename
	newPhoto := &models.Photo{
		UserID:    uid,
		ImageURL:  imageURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	photoID, err := p.storage.SavePhoto(ctx, newPhoto)
	if err != nil {
		p.log.Error("failed to save photo metadata in database", err)
		return 0, fmt.Errorf("failed to save photo metadata in database: %v", err)
	}

	// Возвращаем ID загруженной фотографии
	return photoID, nil
}

// GetPhoto - возвращает данные фотографии по её ID.
// Параметры:
// - ctx: контекст для управления запросом.
// - photoID: ID фотографии для поиска.
// Возвращает данные модели Photo или ошибку.
func (p *Photo) GetPhoto(context context.Context, photoID uint64) (*models.Photo, error) {
	photo, err := p.storage.GetPhoto(context, photoID)
	if err != nil {
		p.log.Error("failed to get photo from storage", err)
		return nil, fmt.Errorf("failed to get photo from storage: %v", err)
	}
	return photo, nil
}

// GetUserPhotos - возвращает список фотографий пользователя по его ID.
// Параметры:
// - ctx: контекст для управления запросом.
// - userID: ID пользователя для поиска его фотографий.
// Возвращает массив моделей Photo или ошибку.
func (p *Photo) GetUserPhotos(context context.Context, userID uint64) ([]*models.Photo, error) {
	photos, err := p.storage.GetUsersPhoto(context, userID)
	if err != nil {
		p.log.Error("failed to get photo from storage", err)
		return nil, fmt.Errorf("failed to get photo from storage: %v", err)
	}
	return photos, nil
}
