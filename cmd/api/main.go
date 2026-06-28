package main

import (
	"context"
	"fmt"
	"github.com/Abhishek48Shah/blog/db"
	"github.com/Abhishek48Shah/blog/internal/config"
	"github.com/Abhishek48Shah/blog/internal/routes"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
)

type Config struct {
	DatabaseUrl string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	cfg, err := config.Load()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	ctx := context.Background()
	pool, err := db.Load(ctx, cfg)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info("database connection established", "port", pool.Config().ConnConfig.Port)
	server := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port), Handler: routes.NewRoutes(pool)}
	slog.Info("server running", "port", server.Addr, "environment", cfg.Environment)
	if err := server.ListenAndServe(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
