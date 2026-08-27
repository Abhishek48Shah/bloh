package config

import (
	"fmt"
	"os"
	"strconv"
)

type Database struct {
	DatabaseURL string
	MaxConn     int
	MinConn     int
}
type Cache struct {
	Adder    string
	Password string
	DB       int
}
type Config struct {
	Port        int
	Environment string
	Db          Database
	Redis       Cache
}

func parseString(s string) (string, error) {
	return s, nil
}
func getenv[T any](key string, defaultValue T, parser func(string) (T, error)) (T, error) {
	if value, ok := os.LookupEnv(key); ok {
		return parser(value)
	}
	return defaultValue, nil
}
func Load() (*Config, error) {
	cfg := Config{}
	var err error
	cfg.Port, err = getenv("PORT", 4000, strconv.Atoi)
	if err != nil {
		return nil, fmt.Errorf("port: %w", err)
	}
	cfg.Environment, _ = getenv("ENVIRONMENT", "devlopment", parseString)
	if err != nil {
		return nil, fmt.Errorf("")
	}
	cfg.Db.DatabaseURL, _ = getenv("DATABASE_URL", "postgres://abhishek:abhishek123@localhost:5432/my_db?sslmode=disable", parseString)
	cfg.Db.MaxConn, err = getenv("MAX_CONNS", 25, strconv.Atoi)
	if err != nil {
		return nil, fmt.Errorf("maxConn: %w", err)
	}
	cfg.Db.MinConn, err = getenv("MIN_CONNS", 5, strconv.Atoi)
	if err != nil {
		return nil, fmt.Errorf("minConn: %w", err)
	}
	cfg.Redis.Adder, _ = getenv("CACHE_URL", "localhost:5432", parseString)
	cfg.Redis.DB, err = getenv("DB", 0, strconv.Atoi)
	if err != nil {
		return nil, err
	}
	cfg.Redis.Password, _ = getenv("PASSWORD", "", parseString)
	return &cfg, nil
}
