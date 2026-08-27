package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Abhishek48Shah/blog/internal/api"
	"github.com/Abhishek48Shah/blog/internal/config"
	"github.com/Abhishek48Shah/blog/internal/database"
	"github.com/joho/godotenv"
)

func Run() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	err := godotenv.Load()
	if err != nil {
		logger.Error("failed to load .env", "error", err.Error())
		os.Exit(1)
	}
	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err.Error())
		os.Exit(1)
	}
	ctx := context.Background()
	pool, err := database.Postgres(ctx, cfg)
	if err != nil {
		logger.Error("failed to load postgres pool", "error", err.Error())
		os.Exit(1)
	}
	cache, err := database.Redis(ctx, cfg)
	if err != nil {
		logger.Error("faile to load redis", "error", err.Error())
		os.Exit(1)
	}
	server := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port), Handler: api.Handler(pool, cache, logger)}
	logger.Info("database connection established", "port", pool.Config().ConnConfig.Port)
	logger.Info("cache connected", "port", 6397)
	logger.Info("server running ")
	err = server.ListenAndServe()
	logger.Error("failed to run server", "error", err)
	defer pool.Close()
	defer cache.Close()
	os.Exit(1)
}
