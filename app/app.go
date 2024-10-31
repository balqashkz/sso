package app

import (
	grpcapp "github.com/balqashkz/sso/app/grpc"
	"github.com/balqashkz/sso/internal/services/auth"
	"github.com/balqashkz/sso/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	// инициализировать хранилище (storage)
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	// init auth service (auth)
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
