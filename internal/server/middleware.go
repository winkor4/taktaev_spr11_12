package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const keyUser ctxKey = "user"

// checkContentType - перехватывает запрос для проверки типа данных
func checkContentType(h http.HandlerFunc, exContentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if !strings.Contains(contentType, exContentType) {
			http.Error(w, "unexpected Content-Type", http.StatusBadRequest)
			return
		}
		h(w, r)
	}
}

// Проверяет авторизацию пользователя
func authorization() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			c, err := r.Cookie("token")
			if err != nil {
				if err == http.ErrNoCookie {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				http.Error(w, "can't read cookie", http.StatusBadRequest)
				return
			}

			tokenStr := c.Value
			claims := new(Claims)

			tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
				return jwtKey, nil
			})

			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				http.Error(w, "can't parse cookie", http.StatusBadRequest)
				return
			}
			if !tkn.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), keyUser, claims.Login)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
