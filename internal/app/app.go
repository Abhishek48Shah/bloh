package app

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Abhishek48Shah/blog/internal/api"
	"github.com/Abhishek48Shah/blog/internal/config"
	"github.com/Abhishek48Shah/blog/internal/database"
	"github.com/Abhishek48Shah/blog/internal/middleware"
	"github.com/Abhishek48Shah/blog/internal/pagination"
	"github.com/Abhishek48Shah/blog/internal/token"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func loadRSAKeys(cfg *config.Config) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	bytPrivKey, err := os.ReadFile(cfg.Jwt.PrivateKeyPath)
	if err != nil {
		return nil, nil, err
	}
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(bytPrivKey)
	if err != nil {
		return nil, nil, err
	}
	bytPubKey, err := os.ReadFile(cfg.Jwt.PublicKeyPath)
	if err != nil {
		return nil, nil, err
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(bytPubKey)
	if err != nil {
		return nil, nil, err
	}
	return privKey, pubKey, nil
}
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
	logger.Info("database connection established", "port", pool.Config().ConnConfig.Port)
	cache, err := database.Redis(ctx, cfg)
	if err != nil {
		logger.Error("faile to load redis", "error", err.Error())
		os.Exit(1)
	}
	logger.Info("cache connected", "port", 6397)
	validator := validator.New()
	privKey, pubKey, err := loadRSAKeys(cfg)
	if err != nil {
		logger.Error("failed to load RSA keys", "error", err.Error())
		os.Exit(1)
	}
	tokenService := token.NewTokenService(privKey, pubKey, cfg.Jwt.Issuer)
	codec := pagination.NewCodec(cfg.HamcSecret)
	middleware := middleware.NewMiddleware(logger, cache, tokenService)
	newRoute := api.NewRoutes(pool, cache, logger, validator, tokenService, middleware, codec)
	server := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port), Handler: newRoute.Handler(), ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError)}
	serverErr := make(chan error, 1)
	go func() {
		logger.Info("server running ", "port", server.Addr)
		serverErr <- server.ListenAndServe()
	}()
	signalCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	select {
	case <-signalCtx.Done():
		logger.Warn("shutdown signal recived")
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "msg", err)
		}
		return
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Warn("gracefull shutdown failed", "error", err)
		if err := server.Close(); err != nil {
			logger.Warn("forced shutdown failed", "error", err)
		}
	}
	logger.Warn("server stopping")
	defer pool.Close()
	defer cache.Close()
}
