package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds every environment-driven setting the app needs.
type Config struct {
	DBURL string

	SessionTimeout   time.Duration
	LockoutThreshold int
	LockoutDuration  time.Duration
}

// Load reads a .env file (if present) and then populates Config from the
// process environment. Missing optional values fall back to sane defaults.
func Load() (*Config, error) {
	// It's fine if .env doesn't exist (e.g. real env vars set some other way,
	// such as inside the Docker container via env_file).
	_ = godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DB_URL is not set")
	}

	sessionTimeout, err := envInt("SESSION_TIMEOUT_MINUTES", 20)
	if err != nil {
		return nil, err
	}

	lockoutThreshold, err := envInt("LOCKOUT_THRESHOLD", 5)
	if err != nil {
		return nil, err
	}

	lockoutDuration, err := envInt("LOCKOUT_DURATION_MINUTES", 15)
	if err != nil {
		return nil, err
	}

	return &Config{
		DBURL:            dbURL,
		SessionTimeout:   time.Duration(sessionTimeout) * time.Minute,
		LockoutThreshold: lockoutThreshold,
		LockoutDuration:  time.Duration(lockoutDuration) * time.Minute,
	}, nil
}

func envInt(key string, def int) (int, error) {
	raw := os.Getenv(key)
	if raw == "" {
		return def, nil
	}
	val, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid value for %s: %w", key, err)
	}
	return val, nil
}
