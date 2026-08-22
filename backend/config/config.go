package config

import (
	"os"

	"github.com/joho/godotenv"
)

type ConfigQW struct {
	QuickwitURL string
}

var AppConfig *ConfigQW

func Load() error {
	_ = godotenv.Load() // ignore error if .env not found

	url, port := os.Getenv("QUICKWIT_URL"), os.Getenv("BACKEND_PORT")
	if url == "" {
		url = "http://quickwit-searcher:7280"
	}
	if port == "" {
		port = "8080"
	}

	AppConfig = &ConfigQW{
		QuickwitURL: url,
	}
	return nil
}
