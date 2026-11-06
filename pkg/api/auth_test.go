package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MaksimKrug/go-praktikum-final-project/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

func TestTokenWithoutExpiration(t *testing.T) {
	password := "test-password"
	secret := []byte(strings.Repeat("s", 32))
	token, err := createToken(password, secret)
	if err != nil {
		t.Fatal(err)
	}
	var claims tokenClaims
	_, err = jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (any, error) {
		return secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		t.Fatal(err)
	}
	if claims.ExpiresAt != nil {
		t.Fatal("token must not have an expiration date")
	}
	if err := validateToken(token, password, secret); err != nil {
		t.Fatalf("token without expiration rejected: %v", err)
	}
	if err := validateToken(token, "changed-password", secret); err == nil {
		t.Fatal("token accepted after password change")
	}
	if err := validateToken(token, password, []byte(strings.Repeat("x", 32))); err == nil {
		t.Fatal("token accepted after secret change")
	}
}

func TestAuthUsesStartupConfig(t *testing.T) {
	t.Setenv("TODO_PORT", "")
	t.Setenv("TODO_PASSWORD", "startup-password")
	t.Setenv("TODO_JWT_SECRET", strings.Repeat("s", 32))
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	a := &apiHandlers{config: cfg}
	token, err := createToken(cfg.Password, cfg.JWTSecret)
	if err != nil {
		t.Fatal(err)
	}

	// Later environment changes must not change the running server's configuration.
	t.Setenv("TODO_PASSWORD", "")
	t.Setenv("TODO_JWT_SECRET", "short")
	called := false
	handler := a.auth(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})
	request := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	w := httptest.NewRecorder()
	handler(w, request)
	if w.Code != http.StatusUnauthorized || called {
		t.Fatal("environment change disabled authentication")
	}
	request.AddCookie(&http.Cookie{Name: "token", Value: token})
	w = httptest.NewRecorder()
	handler(w, request)
	if w.Code != http.StatusNoContent || !called {
		t.Fatal("startup token was rejected after environment change")
	}

	w = httptest.NewRecorder()
	a.signinHandler(w, httptest.NewRequest(http.MethodPost, "/api/signin", strings.NewReader(`{"password":"startup-password"}`)))
	if w.Code != http.StatusOK {
		t.Fatalf("signin stopped using startup configuration: %s", w.Body.String())
	}
}
