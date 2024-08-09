package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/winkor4/taktaev_spr11_12/internal/model"
	"github.com/winkor4/taktaev_spr11_12/internal/users"
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
func authorization(s *Server) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			login, password, ok := r.BasicAuth()
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			auth, err := auth(r.Context(), login, password, "", s.db)
			switch {
			case err != nil:
				http.Error(w, "can't auth", http.StatusInternalServerError)
				return
			case !auth:
				http.Error(w, "can't auth", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), keyUser, login)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func auth(ctx context.Context, l, p, k string, db model.StorageRepo) (bool, error) {
	userManager := users.NewUserManager(l, p, k, db)
	return userManager.CheckAuth(ctx)
}
