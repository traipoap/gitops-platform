package config

import "time"

type JWTConfig struct {
	SecretKey         string
	Expiration        time.Duration
	RefreshExpiration time.Duration
	Issuer            string
}

var Config = JWTConfig{
	SecretKey:         "รักษาความปลอดภัย",
	Expiration:        15 * time.Minute,
	RefreshExpiration: 7 * 24 * time.Hour,
	Issuer:            "High-Performance Centralized Logging & PDPA Compliance System",
}
