package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	DBUrl             string
	JWTSecret         string
	JWTExpireMinutes  int
	TokenLengthBytes  int
	RefreshExpireDays int
	BcryptCost        int
	WebhookURL        string
}

func Load() *Config {
	jwtMin := mustInt(os.Getenv("JWT_EXPIRE_MINUTES"), 15)
	tokBytes := mustInt(os.Getenv("TOKEN_LENGTH_BYTES"), 32)
	refDays := mustInt(os.Getenv("REFRESH_EXPIRE_DAYS"), 7)
	cost := mustInt(os.Getenv("BCRYPT_COST"), 12)
	dburl := os.Getenv("POSTGRES_URL")
	if dburl == "" {
		log.Fatal("POSTGRES_URL is required")
	}
	return &Config{
		DBUrl:             dburl,
		JWTSecret:         os.Getenv("JWT_SECRET"),
		JWTExpireMinutes:  jwtMin,
		TokenLengthBytes:  tokBytes,
		RefreshExpireDays: refDays,
		BcryptCost:        cost,
		WebhookURL:        os.Getenv("WEBHOOK_URL"),
	}
}

func mustInt(val string, def int) int {
	if i, err := strconv.Atoi(val); err == nil {
		return i
	}
	return def
}
