package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Seed     SeedConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret  string
	Expires time.Duration
}

type SeedConfig struct {
	AdminUsername string
	AdminPassword string
	AdminRealName string
}

func Load() Config {
	_ = godotenv.Load()

	mode := env("GIN_MODE", "debug")
	gin.SetMode(mode)

	return Config{
		Server: ServerConfig{
			Port: env("SERVER_PORT", "8989"),
			Mode: mode,
		},
		Database: DatabaseConfig{
			Host:     env("DB_HOST", "127.0.0.1"),
			Port:     env("DB_PORT", "3306"),
			User:     env("DB_USER", "root"),
			Password: env("DB_PASSWORD", "root"),
			Name:     env("DB_NAME", "cy11_dsphk"),
		},
		JWT: JWTConfig{
			Secret:  env("JWT_SECRET", "cy11-dsphk-local-dev-secret"),
			Expires: time.Duration(envInt("JWT_EXPIRES_HOURS", 168)) * time.Hour,
		},
		Seed: SeedConfig{
			AdminUsername: env("ADMIN_USERNAME", "admin"),
			AdminPassword: env("ADMIN_PASSWORD", "admin123456"),
			AdminRealName: env("ADMIN_REAL_NAME", "超级管理员"),
		},
	}
}

func (d DatabaseConfig) DSN(databaseName string) string {
	if databaseName == "" {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", d.User, d.Password, d.Host, d.Port)
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", d.User, d.Password, d.Host, d.Port, databaseName)
}

func env(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
