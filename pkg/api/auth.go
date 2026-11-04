package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type tokenClaims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

// ValidateAuthConfig checks the signing secret when authentication is enabled.
func ValidateAuthConfig() error {
	if os.Getenv("TODO_PASSWORD") == "" {
		return nil
	}
	_, err := jwtSecret()
	return err
}

func jwtSecret() ([]byte, error) {
	secret := os.Getenv("TODO_JWT_SECRET")
	if len(secret) < 32 {
		return nil, fmt.Errorf("TODO_JWT_SECRET must contain at least 32 bytes")
	}
	return []byte(secret), nil
}

func passwordHash(password string, secret []byte) string {
	hash := hmac.New(sha256.New, secret)
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func createToken(password string) (string, error) {
	secret, err := jwtSecret()
	if err != nil {
		return "", err
	}

	claims := tokenClaims{
		PasswordHash: passwordHash(password, secret),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func validateToken(tokenString, password string) error {
	secret, err := jwtSecret()
	if err != nil {
		return err
	}

	var claims tokenClaims

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (any, error) {
			return secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return fmt.Errorf("parse token: %w", err)
	}
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	expectedHash := passwordHash(password, secret)
	if !hmac.Equal([]byte(claims.PasswordHash), []byte(expectedHash)) {
		return fmt.Errorf("password has changed")
	}

	return nil
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("TODO_PASSWORD")

		if password != "" {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			if err := validateToken(cookie.Value, password); err != nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}

		next(w, r)
	}
}
