package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DBHost           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBPort           string
	JWTSecret        string
	CORSAllowOrigins []string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:           os.Getenv("POSTGRES_HOSTNAME"),
		DBUser:           os.Getenv("POSTGRES_USER"),
		DBPassword:       os.Getenv("POSTGRES_PASSWORD"),
		DBName:           os.Getenv("POSTGRES_DB"),
		DBPort:           os.Getenv("POSTGRES_PORT"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		CORSAllowOrigins: strings.Split(os.Getenv("CORS_ALLOW_ORIGINS"), ","),
	}
}

func (c *Config) GetDBConnString() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort,
	)
}
