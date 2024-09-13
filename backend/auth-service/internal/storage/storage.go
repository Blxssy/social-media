package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/Blxssy/social-media/auth-service/internal/config"
	"github.com/Blxssy/social-media/auth-service/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
)

type Storage interface {
	SaveUser(ctx context.Context, username string, email string, passHash []byte) error
	User(ct context.Context, email string) (*models.User, error)
	IsAdmin(ct context.Context, userID int) (bool, error)
}

type storage struct {
	db *gorm.DB
}

func NewStorage(logger *slog.Logger, config *config.Config) Storage {
	db, err := connectDatabase(config)
	if err != nil {
		logger.Error("Failure database connection")
		panic(err)
	}
	logger.Info("Successfully connection to database")

	db.AutoMigrate(&models.User{})

	passHash, err := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	Admin := models.User{
		Model:    gorm.Model{},
		Username: "Anton Sonin",
		Email:    "test@test.com",
		PassHash: string(passHash),
		IsAdmin:  true,
	}
	db.Create(&Admin)

	return &storage{
		db: db,
	}
}

func connectDatabase(config *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.Username,
		config.Database.Name, config.Database.Password)
	return gorm.Open(postgres.Open(dsn))
}

func (s *storage) SaveUser(ctx context.Context, username string, email string, passHash []byte) error {
	u, _ := s.findByEmail(ctx, email)
	if u != nil {
		return errors.New("User already exists")
	}

	user := &models.User{
		Username: username,
		Email:    email,
		PassHash: string(passHash),
	}

	err := s.db.Create(user).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) User(ctx context.Context, email string) (*models.User, error) {
	user, err := s.findByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *storage) IsAdmin(ctx context.Context, userID int) (bool, error) {
	user, err := s.findByID(ctx, userID)
	if err != nil {
		return false, err
	}

	if user == nil {
		return false, nil
	}

	return user.IsAdmin, nil
}

func (s *storage) findByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *storage) findByID(ctx context.Context, uid int) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ?", uid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
