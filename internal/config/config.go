package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string
	Port        string
	DatabaseUrl string

	Auth struct {
		JwtAccessSecret    string
		JwtRefreshSecret   string
		AccessTokenMinutes int
		RefreshTokenDays   int
	}

	CorsOrigin string

	Cloudinary struct {
		CloudName string
		ApiKey    string
		ApiSecret string
	}

	Smtp struct {
		Host     string
		Port     string
		Username string
		Password string
		From     string
	}

	MaxUploadMb int
}

func MustLoad() Config {
	_ = godotenv.Load()
	var cfg Config
	cfg.Env = getEnv("APP_ENV")
	cfg.Port = getEnvOrDefault("PORT", "5000")
	cfg.DatabaseUrl = getEnv("DATABASE_URL")
	cfg.Auth.JwtAccessSecret = getEnv("JWT_ACCESS_SECRET")
	cfg.Auth.JwtRefreshSecret = getEnv("JWT_REFRESH_SECRET")
	cfg.Auth.AccessTokenMinutes = getEnvOrDefaultInt("ACCESS_TOKEN_MINUTES", "15")
	cfg.Auth.RefreshTokenDays = getEnvOrDefaultInt("REFRESH_TOKEN_DAYS", "7")
	cfg.CorsOrigin = getEnvOrDefault("CORS_ORIGIN", "http://localhost:3000")
	cfg.Cloudinary.CloudName = getEnv("CLOUDINARY_CLOUD_NAME")
	cfg.Cloudinary.ApiKey = getEnv("CLOUDINARY_API_KEY")
	cfg.Cloudinary.ApiSecret = getEnv("CLOUDINARY_API_SECRET")
	cfg.Smtp.Host = getEnv("SMTP_HOST")
	cfg.Smtp.Port = getEnv("SMTP_PORT")
	cfg.Smtp.Username = getEnv("SMTP_USERNAME")
	cfg.Smtp.Password = getEnv("SMTP_PASSWORD")
	cfg.Smtp.From = getEnv("SMTP_FROM")
	cfg.MaxUploadMb = getEnvOrDefaultInt("MAX_UPLOAD_MB", "5")

	return cfg
}

func getEnv(key string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		panic(fmt.Sprintf("%v is not set", key))
	}
	return v
}

func getEnvOrDefault(key string, defaultValue string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return defaultValue
	}
	return v
}

func getEnvOrDefaultInt(key string, defaultValue string) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		d, err := strconv.Atoi(defaultValue)
		if err != nil {
			panic(fmt.Sprintf("invalid default value for %v: %v", key, defaultValue))
		}
		return d
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		panic(fmt.Sprintf("invalid value for %v: %v", key, v))
	}
	return i
}

// func getEnvOrDefaultBool(key string, defaultValue string) bool {
// 	v := strings.TrimSpace(os.Getenv(key))
// 	if v == "" {
// 		d, err := strconv.ParseBool(defaultValue)
// 		if err != nil {
// 			panic(fmt.Sprintf("invalid default value for %v: %v", key, defaultValue))
// 		}
// 		return d
// 	}

// 	b, err := strconv.ParseBool(v)
// 	if err != nil {
// 		panic(fmt.Sprintf("invalid value for %v: %v", key, v))
// 	}
// 	return b
// }

// func getEnvOrDefaultDuration(key string, defaultValue string) time.Duration {
// 	v := strings.TrimSpace(os.Getenv(key))
// 	if v == "" {
// 		d, err := time.ParseDuration(defaultValue)
// 		if err != nil {
// 			panic(fmt.Sprintf("invalid default value for %v: %v", key, defaultValue))
// 		}
// 		return d
// 	}

// 	d, err := time.ParseDuration(v)
// 	if err != nil {
// 		panic(fmt.Sprintf("invalid value for %v: %v", key, v))
// 	}
// 	return d
// }
