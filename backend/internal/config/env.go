package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	FrontendURL          string
	DBHost               string
	DBDatabase           string
	DBUsername           string
	DBPassword           string
	DBPort               string
	SeedAdminEmail       string
	SeedAdminPassword    string
	SeedAdminFirstName   string
	SeedAdminMiddleName  string
	SeedAdminLastName    string
	AppName              string
	PasetoSymmetricKey   string
	AccessTokenDuration  string
	RefreshTokenDuration string
	CSRFTokenDuration    string
	InviteDuration       string
	PasswordResetDuration string
	GoEnv                string
	LogLevel             string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppName:              getEnv("APP_NAME", "secure-web-helpdesk"),
		Port:                 getEnv("PORT", "8080"),
		FrontendURL:          getEnv("FRONTEND_URL", "http://localhost:5173"),
		GoEnv:                getEnv("GO_ENV", "development"),
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBDatabase:           getEnv("DB_DATABASE", "helpdesk"),
		DBUsername:           getEnv("DB_USERNAME", "admin"),
		DBPassword:           getEnv("DB_PASSWORD", "password"),
		DBPort:               getEnv("DB_PORT", "3306"),
		SeedAdminEmail:       getEnv("SEED_ADMIN_EMAIL", "admin@helpdesk.ie"),
		SeedAdminPassword:    getEnv("SEED_ADMIN_PASSWORD", "password"),
		SeedAdminFirstName:   getEnv("SEED_ADMIN_FIRST_NAME", "cyber"),
		SeedAdminMiddleName:  getEnv("SEED_ADMIN_MIDDLE_NAME", ""),
		SeedAdminLastName:    getEnv("SEED_ADMIN_LAST_NAME", "security"),
		PasetoSymmetricKey:   getEnv("PASETO_SYMMETRIC_KEY", "12345678901234567890123456789012"),
		AccessTokenDuration:  getEnv("ACCESS_TOKEN_DURATION", "15m"),
		RefreshTokenDuration: getEnv("REFRESH_TOKEN_DURATION", "168h"),
		CSRFTokenDuration:    getEnv("CSRF_TOKEN_DURATION", "60m"),
		InviteDuration:        getEnv("INVITE_TTL", "72h"),
		PasswordResetDuration: getEnv("PASSWORD_RESET_TTL", "1h"),
		LogLevel:              getEnv("LOG_LEVEL", "info"),
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

func (c Config) AccessTokenTTL() time.Duration {
	return parseDuration(c.AccessTokenDuration, 15*time.Minute)
}

func (c Config) RefreshTokenTTL() time.Duration {
	return parseDuration(c.RefreshTokenDuration, 7*24*time.Hour)
}

func (c Config) CSRFTTL() time.Duration {
	return parseDuration(c.CSRFTokenDuration, 60*time.Minute)
}

func (c Config) InviteTTL() time.Duration {
	return parseDuration(c.InviteDuration, 72*time.Hour)
}

func (c Config) PasswordResetTTL() time.Duration {
	return parseDuration(c.PasswordResetDuration, time.Hour)
}

func (c Config) TokenIssuer() string {
	name := strings.TrimSpace(c.AppName)
	if name == "" {
		return "secure-web-helpdesk"
	}
	return name
}

func (c Config) TokenAudience() string {
	return fmt.Sprintf("%s-web", c.TokenIssuer())
}

func parseDuration(raw string, fallback time.Duration) time.Duration {
	parsed, err := time.ParseDuration(strings.TrimSpace(raw))
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
