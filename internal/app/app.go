package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/vindosVP/snprofiles/cmd/config"
	"github.com/vindosVP/snprofiles/internal/app/grpc"
	"github.com/vindosVP/snprofiles/internal/storage"
)

type App struct {
	GRPCServer *grpc.App
}

func New(log zerolog.Logger, cfg *config.Config) *App {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, postgresConn(cfg))
	if err != nil {
		panic(fmt.Errorf("could not connect to postgres: %w", err))
	}
	p := storage.NewPostgresStorage(pool)
	grpcApp := grpc.New(log, p, cfg.GRPC.Port)
	return &App{
		GRPCServer: grpcApp,
	}
}

func postgresConn(cfg *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Database,
	)
}
