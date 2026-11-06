package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type tokenClaims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

func passwordHash(password string, secret []byte) string {
	hash := hmac.New(sha256.New, secret)
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func createToken(password string, secret []byte) (string, error) {
	claims := tokenClaims{
		PasswordHash: passwordHash(password, secret),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func validateToken(tokenString, password string, secret []byte) error {
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

func (a *apiHandlers) auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		password := a.config.Password

		if password != "" {
			cookie, err := r.Cookie("token")
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "authentication required"})
				return
			}

			if err := validateToken(cookie.Value, password, a.config.JWTSecret); err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "authentication required"})
				return
			}
		}

		next(w, r)
	}
}
