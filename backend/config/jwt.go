package config

import (
	"errors"
	"os"
	"strings"
	"time"
)

type JWTConfig struct {
	SecretKey         string
	Expiration        time.Duration
	RefreshExpiration time.Duration
	Issuer            string
}

// Config holds JWT settings. SecretKey is populated by Load() from the
// JWT_SECRET environment variable — it must never be hardcoded.
var Config = JWTConfig{
	Expiration:        15 * time.Minute,
	RefreshExpiration: 7 * 24 * time.Hour,
	Issuer:            "High-Performance Centralized Logging",
}

// loadJWT validates JWT settings and sets Config.SecretKey from the environment.
func loadJWT() error {
	secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if secret == "" {
		return errors.New("JWT_SECRET environment variable is required")
	}
	Config.SecretKey = secret
	return nil
}
