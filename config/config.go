package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	MongoURI  string
	MongoDB   string
	JwtSecret string
	JwtTtl    time.Duration
}

func Load() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	ttl, err := parseDuration(env("JWT_TTL", "24h"))
	if err != nil {
		return Config{}, fmt.Errorf("JWT_TTL: %w", err)
	}

	cfg := Config{
		Port:      env("PORT", ":8080"),
		MongoURI:  env("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:   env("MONGO_DB", "user_management"),
		JwtSecret: os.Getenv("JWT_SECRET"),
		JwtTtl:    ttl,
	}

	if cfg.JwtSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}
	return cfg, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseDuration(v string) (time.Duration, error) {
	return time.ParseDuration(v)
}
