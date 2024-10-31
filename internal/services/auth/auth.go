package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/balqashkz/sso/internal/domain/models"
	"github.com/balqashkz/sso/internal/lib/jwt"
	"github.com/balqashkz/sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, username string, hashPass []byte) (int64, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

/*
Login checks if user with given credentials exists in the system and returns token

if user exists, but password is incorrect, returns error
if user doesn't exist, returns error
*/
func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	const op = "auth.Login"

	log := slog.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("attempting to login")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", slog.String("error", err.Error()))

			return "", fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		a.log.Error("failed to get user", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("Invalid credentials", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged is successfully")
	return token, nil
}

/*
RegisterNewUser registers new user in the system and returns user ID
if user with given username already exists, returns error
*/
func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", slog.String("error", err.Error()))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	userID, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", slog.String("error", err.Error()))

			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		log.Error("failed to save user", slog.String("error", err.Error()))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered", slog.String("email", email))

	return userID, nil
}
