package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port      string
	DBFile    string
	Password  string
	JWTSecret []byte
}

// Load reads and validates the environment once, before the server starts.
func Load() (Config, error) {
	cfg := Config{
		Port:      os.Getenv("TODO_PORT"),
		DBFile:    os.Getenv("TODO_DBFILE"),
		Password:  os.Getenv("TODO_PASSWORD"),
		JWTSecret: []byte(os.Getenv("TODO_JWT_SECRET")),
	}
	if cfg.Port == "" {
		cfg.Port = "7540"
	}
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		return Config{}, fmt.Errorf("port %q is not a number: %w", cfg.Port, err)
	}
	if port < 1 || port > 65535 {
		return Config{}, fmt.Errorf("port is outside [1, 65535]: %s", cfg.Port)
	}
	if cfg.DBFile == "" {
		cfg.DBFile = "scheduler.db"
	}
	if cfg.Password != "" && len(cfg.JWTSecret) < 32 {
		return Config{}, fmt.Errorf("jwt secret must contain at least 32 bytes (TODO_JWT_SECRET)")
	}
	return cfg, nil
}
