package config

import (
	"fmt"
	"log"
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
		PublicHost:           getEnv("PUBLIC_HOST"),
		Port:                 getEnv("PORT"),
		DBUser:               getEnv("DB_USER"),
		DBPassword:           getEnv("DB_PASSWORD"),
		DBAddress:            fmt.Sprintf("%s:%s", getEnv("DB_HOST"), getEnv("DB_PORT")),
		DBName:               getEnv("DB_NAME"),
		JWTAccessExpiration:  getEnvAsInt("JWT_ACCESS_EXPIRATION"),
		JWTRefreshExpiration: getEnvAsInt("JWT_REFRESH_EXPIRATION"),
		JWTSecret:            getEnv("JWT_SECRET"),
	}
}
func getEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("environment variable %s not set", key)
	}
	return value
}

func getEnvAsInt(key string) int64 {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("environment variable %s not set", key)
	}
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Fatalf("environment variable %s is not a valid integer", key)
	}
	return i
}
