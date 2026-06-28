package config

import (
	"os"
	"strconv"
)

type Db struct {
	Dsn string
}
type Config struct {
	Port        int
	Environment string
	DB          Db
}

func stringParser(s string) (string, error) {
	return s, nil
}

func getEnv[T any](key string, defaultValue T, parse func(string) (T, error)) (T, error) {
	if value, ok := os.LookupEnv(key); ok {
		return parse(value)
	}
	return defaultValue, nil
}
func Load() (*Config, error) {
	port, err := getEnv("PORT", 4000, strconv.Atoi)
	if err != nil {
		return nil, err
	}
	env, err := getEnv("ENVIRONMENT", "devlopment", stringParser)
	if err != nil {
		return nil, err
	}
	dsn, err := getEnv("DATABASE_URL", "postgres://abhishek:abhsihek123@localhost:5432/my_db?sslmode=disable", stringParser)
	if err != nil {
		return nil, err
	}
	return &Config{Port: port, Environment: env, DB: Db{Dsn: dsn}}, nil
}
