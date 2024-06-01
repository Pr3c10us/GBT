package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type PGDB struct {
	Username string
	Password string
	Host     string
	Port     int
	Name     string
}

type EnvironmentVariables struct {
	Port                  string
	JWTSecret             string
	CookieSecret          string
	SessionSecret         string
	SessionMaxAge         int
	ProductionEnvironment bool
	AuthRedirectUrl       string
	PGDB                  *PGDB
	ClientDomain          string
}

func LoadEnvironment() *EnvironmentVariables {
	err := godotenv.Load("C:/Users/owopr/Documents/go/GBT/.dev.env")
	if err != nil {

		log.Fatal(err)
	}
	return &EnvironmentVariables{
		Port:                  getEnv("PORT", ":5000"),
		JWTSecret:             getEnvOrError("JWT_SECRET"),
		CookieSecret:          getEnvOrError("COOKIE_SECRET"),
		SessionSecret:         getEnvOrError("SESSIONS_SECRET"),
		SessionMaxAge:         getEnvAsInt("SESSION_MAX_AGE", 86400*300),
		ProductionEnvironment: getEnvAsBool("PRODUCTION_ENVIRONMENT", false),
		PGDB: &PGDB{
			Username: getEnv("PG_DB_USERNAME", "postgres"),
			Password: getEnvOrError("PG_DB_PASSWORD"),
			Host:     getEnv("PG_DB_HOST", "127.0.0.1"),
			Port:     getEnvAsInt("PG_DB_PORT", 5432),
			Name:     getEnvOrError("PG_DB_NAME"),
		},
		ClientDomain: getEnv("CLIENT_DOMAIN", "localhost"),
	}
}

func getEnvOrError(key string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	panic("Environment variable not set")
}

func getEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	value, exist := os.LookupEnv(key)
	if exist {
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			log.Panicf("Environment variable \"%v\" not set properly", key)
		}
		return valueInt
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	value, exist := os.LookupEnv(key)
	if exist {
		valueBool, err := strconv.ParseBool(value)
		if err != nil {
			log.Panicf("Environment variable \"%v\" not set properly", key)
		}
		return valueBool
	}
	return fallback
}
