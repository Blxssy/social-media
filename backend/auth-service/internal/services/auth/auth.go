package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/Blxssy/social-media/auth-service/internal/models"
	"github.com/Blxssy/social-media/auth-service/pkg/token"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	tokenSaver  TokenSaver
}

type UserSaver interface {
	SaveUser(ctx context.Context, username string, email string, passHash []byte) error
}

type UserProvider interface {
	User(ct context.Context, email string) (*models.User, error)
	IsAdmin(ct context.Context, userID int) (bool, error)
}

type TokenSaver interface {
	SaveTokens(ctx context.Context, uid uint, accessToken string, refreshToken string) error
}

func New(
	log *slog.Logger,
	usrSaver UserSaver,
	usrProvider UserProvider,
	tokenSaver TokenSaver,
) *Auth {
	return &Auth{
		log:         log,
		usrSaver:    usrSaver,
		usrProvider: usrProvider,
		tokenSaver:  tokenSaver,
	}
}

func (a *Auth) Register(ctx context.Context, username, email, password string) (string, string, error) {
	const op = "auth.Register"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	//log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash")

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	user := models.User{
		Username: username,
		Email:    email,
		PassHash: string(passHash),
	}

	err = a.usrSaver.SaveUser(ctx, username, user.Email, passHash)
	if err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err := token.GetNewTokens(user.ID)
	if err != nil {
		log.Error("failed to get new tokens")
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	a.tokenSaver.SaveTokens(ctx, user.ID, accessToken, refreshToken)

	return accessToken, refreshToken, nil
}

func (a *Auth) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		return "", "", err
	}

	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	err = bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password))
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	accessToken, refreshToken, err := token.GetNewTokens(user.ID)
	if err != nil {
		log.Error("failed to get new tokens")
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	a.tokenSaver.SaveTokens(ctx, user.ID, accessToken, refreshToken)

	return accessToken, refreshToken, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int) (bool, error) {
	isAdmin, err := a.usrProvider.IsAdmin(ctx, userID)
	if err != nil {
		return false, err
	}

	return isAdmin, nil
}
