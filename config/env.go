package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost           string
	Port                 string
	DBUser               string
	DBPassword           string
	DBAddress            string
	DBName               string
	JWTSecret            string
	JWTAccessExpiration  int64
	JWTRefreshExpiration int64
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		PublicHost:           getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                 getEnv("PORT", "8080"),
		DBUser:               getEnv("DB_USER", "root"),
		DBPassword:           getEnv("DB_PASSWORD", "admin"),
		DBAddress:            fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:               getEnv("DB_NAME", "gotask"),
		JWTAccessExpiration:  getEnvAsInt("JWT_Access_EXPIRATION", 3600),
		JWTRefreshExpiration: getEnvAsInt("JWT_REFRESH_EXPIRATION", 3600*24*7),
		JWTSecret:            getEnv("JWT_SECRET)", "secret"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return i
		}
	}

	return fallback
}
