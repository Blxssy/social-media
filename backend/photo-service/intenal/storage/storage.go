package storage

import (
	"context"
	"fmt"
	"github.com/Blxssy/social-media/photo-service/intenal/config"
	"github.com/Blxssy/social-media/photo-service/intenal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
)

// Storage - интерфейс, описывающий методы для работы с фотографиями в базе данных.
// Методы:
// - SavePhoto: сохраняет фотографию в базе данных.
// - GetPhoto: получает фотографию по её ID.
// - GetUsersPhoto: получает все фотографии конкретного пользователя по его ID.
type Storage interface {
	SavePhoto(ctx context.Context, photo *models.Photo) (uint64, error)
	GetPhoto(ctx context.Context, id uint64) (*models.Photo, error)
	GetUsersPhoto(ctx context.Context, uid uint64) ([]*models.Photo, error)
}

// storage - структура, реализующая интерфейс Storage.
// Поля:
// - db: объект GORM для работы с базой данных.
// - logger: логгер для записи событий.
type storage struct {
	db     *gorm.DB
	logger *slog.Logger
}

// NewStorage - функция для создания нового экземпляра хранилища.
// Параметры:
// - cfg: конфигурация, содержащая параметры подключения к базе данных.
// - logger: логгер для записи событий.
// Возвращает экземпляр интерфейса Storage или паникует в случае ошибки соединения.
func NewStorage(cfg *config.Config, logger *slog.Logger) Storage {
	db, err := connectDatabase(cfg)
	if err != nil {
		logger.Error("Failure database connection")
		panic(err)
	}
	logger.Info("Successfully connection to database")

	return &storage{
		db:     db,
		logger: logger,
	}
}

// connectDatabase - функция для подключения к базе данных.
// Параметры:
// - config: конфигурация, содержащая параметры подключения.
// Возвращает объект GORM для работы с базой данных и ошибку, если она произошла.
func connectDatabase(config *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.Username,
		config.Database.Name, config.Database.Password)
	return gorm.Open(postgres.Open(dsn))
}

// SavePhoto - метод для сохранения фотографии в базе данных.
// Параметры:
// - ctx: контекст для управления запросом.
// - photo: указатель на модель Photo, которую нужно сохранить.
// Возвращает ID сохранённой фотографии или ошибку.
func (s *storage) SavePhoto(ctx context.Context, photo *models.Photo) (uint64, error) {
	if err := s.db.FirstOrCreate(photo).Error; err != nil {
		return 0, err
	}
	return photo.ID, nil
}

// GetPhoto - метод для получения фотографии по её ID.
// Параметры:
// - ctx: контекст для управления запросом.
// - id: ID фотографии для поиска.
// Возвращает указатель на модель Photo или ошибку.
func (s *storage) GetPhoto(ctx context.Context, id uint64) (*models.Photo, error) {
	var photo models.Photo
	if err := s.db.First(&photo, id).Error; err != nil {
		return nil, err
	}
	return &photo, nil
}

// GetUsersPhoto - метод для получения всех фотографий пользователя по его ID.
// Параметры:
// - ctx: контекст для управления запросом.
// - uid: ID пользователя, для которого нужно получить фотографии.
// Возвращает массив указателей на модели Photo или ошибку.
func (s *storage) GetUsersPhoto(ctx context.Context, uid uint64) ([]*models.Photo, error) {
	var photos []*models.Photo
	if err := s.db.Where("user_id = ?", uid).Find(&photos).Error; err != nil {
		return nil, err
	}
	return photos, nil
}
