package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	DBHost              string
	DBDatabase          string
	DBUsername          string
	DBPassword          string
	DBPort              string
	SeedAdminEmail      string
	SeedAdminPassword   string
	SeedAdminFirstName  string
	SeedAdminMiddleName string
	SeedAdminLastName   string
	GoEnv               string
	LogLevel            string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		Port:                getEnv("PORT", "8080"),
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBDatabase:          getEnv("DB_DATABASE", "helpdesk"),
		DBUsername:          getEnv("DB_USERNAME", "admin"),
		DBPassword:          getEnv("DB_PASSWORD", "password"),
		DBPort:              getEnv("DB_PORT", "3306"),
		SeedAdminEmail:      getEnv("SEED_ADMIN_EMAIL", "x25166280@Student.ncirl.ie"),
		SeedAdminPassword:   getEnv("SEED_ADMIN_PASSWORD", "password"),
		SeedAdminFirstName:  getEnv("SEED_ADMIN_FIRST_NAME", "cyber"),
		SeedAdminMiddleName: getEnv("SEED_ADMIN_MIDDLE_NAME", ""),
		SeedAdminLastName:   getEnv("SEED_ADMIN_LAST_NAME", "security"),
		GoEnv:               getEnv("GO_ENV", "development"),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
	}
}

func (c Config) MySQLDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.DBUsername,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBDatabase,
	)
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
