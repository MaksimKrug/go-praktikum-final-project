package config

import (
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	for _, tc := range []struct {
		name, port, password, secret string
		wantError                    bool
	}{
		{name: "defaults"},
		{name: "password without secret", password: "test", wantError: true},
		{name: "short secret", password: "test", secret: "short", wantError: true},
		{name: "auth enabled", port: "8080", password: "test", secret: strings.Repeat("s", 32)},
		{name: "invalid port", port: "abc", wantError: true},
		{name: "zero port", port: "0", wantError: true},
		{name: "large port", port: "65536", wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("TODO_PORT", tc.port)
			t.Setenv("TODO_DBFILE", "")
			t.Setenv("TODO_PASSWORD", tc.password)
			t.Setenv("TODO_JWT_SECRET", tc.secret)
			t.Setenv("TODO_CORS_ORIGINS", "")
			cfg, err := Load()
			if (err != nil) != tc.wantError {
				t.Fatalf("Load error = %v, wantError = %v", err, tc.wantError)
			}
			if err != nil {
				return
			}
			wantPort := tc.port
			if wantPort == "" {
				wantPort = "7540"
			}
			if cfg.Port != wantPort || cfg.DBFile != "scheduler.db" || cfg.Password != tc.password || string(cfg.JWTSecret) != tc.secret {
				t.Fatal("configuration does not match environment and defaults")
			}
		})
	}
}

func TestLoadCORSOrigins(t *testing.T) {
	t.Setenv("TODO_PORT", "8080")
	t.Setenv("TODO_PASSWORD", "")
	for _, tc := range []struct {
		value     string
		want      string
		wantError bool
	}{
		{"", "http://localhost:8080,http://127.0.0.1:8080", false},
		{"http://localhost:3000, https://app.example", "http://localhost:3000,https://app.example", false},
		{"*", "", true},
		{"http://localhost:3000,", "", true},
		{"https://app.example/path", "", true},
		{"https://app.example/", "", true},
		{"ftp://app.example", "", true},
	} {
		t.Run(tc.value, func(t *testing.T) {
			t.Setenv("TODO_CORS_ORIGINS", tc.value)
			cfg, err := Load()
			if (err != nil) != tc.wantError {
				t.Fatalf("Load error = %v, wantError = %v", err, tc.wantError)
			}
			if err == nil && strings.Join(cfg.CORSAllowedOrigins, ",") != tc.want {
				t.Fatalf("unexpected origins: %v", cfg.CORSAllowedOrigins)
			}
		})
	}
}
