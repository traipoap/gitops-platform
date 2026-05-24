package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type ConfigQW struct {
	QuickwitURL string
}

var AppConfig *ConfigQW

func Load() error {
	_ = godotenv.Load() // ignore error if .env not found

	url := os.Getenv("QUICKWIT_URL")
	if url == "" {
		return fmt.Errorf("QUICKWIT_URL environment variable is required")
	}

	AppConfig = &ConfigQW{
		QuickwitURL: url,
	}
	return nil
}
