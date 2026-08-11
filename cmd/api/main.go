package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Abhishek48Shah/blog/db"
	"github.com/Abhishek48Shah/blog/internal/config"
	"github.com/Abhishek48Shah/blog/internal/routes"
	"github.com/joho/godotenv"
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
	postgresPool, err := db.LoadPostgres(ctx, cfg)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	cache, err := db.LoadRedis(ctx)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info("database connection established", "port", postgresPool.Config().ConnConfig.Port)
	slog.Info("cache connection established", "port", 6379)
	server := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port), Handler: routes.NewRoutes(postgresPool, cache)}
	slog.Info("server running", "port", server.Addr, "environment", cfg.Environment)
	if err := server.ListenAndServe(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	fmt.Println("hello world")
}
