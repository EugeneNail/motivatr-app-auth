package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	App App
	Db  Db
}

type App struct {
	Name string
}

type Db struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func New() (*Config, error) {
	dbPort, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		return nil, fmt.Errorf("converting a database port from string to integer: %w", err)
	}

	return &Config{
		App: App{
			Name: os.Getenv("APP_NAME"),
		},
		Db: Db{
			Host:     os.Getenv("DB_HOST"),
			Name:     os.Getenv("POSTGRES_DB"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Port:     dbPort,
		},
	}, nil
}
