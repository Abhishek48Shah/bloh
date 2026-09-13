package config

import (
	"fmt"
	"os"
	"strconv"
)

type Database struct {
	User     string
	Host     string
	Port     int
	Password string
	DbName   string
	MaxConn  int
	MinConn  int
	SSLMode  string
}
type Cache struct {
	Port     int
	Host     string
	DbName   int
	Password string
}
type JWTConfig struct {
	PrivateKeyPath string
	PublicKeyPath  string
	Issuer         string
}
type Config struct {
	Port        int
	Environment string
	Db          Database
	Redis       Cache
	Jwt         JWTConfig
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
	cfg.Db.User, _ = getenv("POSTGRES_USER", "abhishek", parseString)
	cfg.Db.Password, _ = getenv("POSTGRES_PASSWORD", "abhishek123", parseString)
	cfg.Db.Host, _ = getenv("POSTGRES_HOST", "localhost", parseString)
	cfg.Db.DbName, _ = getenv("POSTGRES_DATABASE", "my_db", parseString)
	cfg.Db.SSLMode, _ = getenv("POSTGRES_SSL_MODE", "disable", parseString)

	cfg.Db.Port, err = getenv("POSTGRES_PORT", 5432, strconv.Atoi)
	if err != nil {
		return nil, fmt.Errorf("postgres port: %w", err)
	}
	cfg.Environment, _ = getenv("ENVIRONMENT", "devlopment", parseString)
	cfg.Db.MaxConn, err = getenv("MAX_CONNS", 25, strconv.Atoi)
	if err != nil {
		return nil, fmt.Errorf("maxConn: %w", err)
	}
	cfg.Db.MinConn, err = getenv("MIN_CONNS", 5, strconv.Atoi)
	if err != nil {
		return nil, fmt.Errorf("minConn: %w", err)
	}
	cfg.Redis.Host, _ = getenv("REDIS_HOST", "localhost", parseString)
	cfg.Redis.Password, _ = getenv("REDIS_PASSWORD", "", parseString)
	cfg.Redis.Port, err = getenv("REDIS_PORT", 6379, strconv.Atoi)
	if err != nil {
		return nil, fmt.Errorf("redis port: %w", err)
	}
	cfg.Redis.DbName, err = getenv("REDIS_DATABASE", 0, strconv.Atoi)
	if err != nil {
		return nil, fmt.Errorf("redis db: %w", err)
	}
	cfg.Jwt.PrivateKeyPath, _ = getenv("PRIVATE_KEY_PATH", "/home/abhishek/buffer/blog/keys/private.pem", parseString)
	cfg.Jwt.PublicKeyPath, _ = getenv("PUBLIC_KEY_PATH", "/home/abhishek/buffer/blog/keys/public.pem", parseString)
	cfg.Jwt.Issuer, _ = getenv("ISSUER", "https://mydomain.com", parseString)
	return &cfg, nil
}
