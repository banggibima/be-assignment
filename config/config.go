package config

import (
	"os"
)

type HTTP struct {
	Port string
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

type Config struct {
	HTTP     HTTP
	Postgres Postgres
}

func Load() (*Config, error) {
	config := &Config{
		HTTP: HTTP{
			Port: os.Getenv("HTTP_PORT"),
		},
		Postgres: Postgres{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: os.Getenv("POSTGRES_DATABASE"),
			SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
		},
	}

	return config, nil
}
