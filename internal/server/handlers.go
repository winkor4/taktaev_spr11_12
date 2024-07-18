package server

import (
	"encoding/json"
	"net/http"

	"github.com/winkor4/taktaev_spr11_12/internal/model"
)

// Регистрация пользователя в приложении
func register(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var parameters model.RegisterRequest
		err := json.NewDecoder(r.Body).Decode(&parameters)
		if err != nil {
			http.Error(w, "Can't read body", http.StatusBadRequest)
			return
		}

		if parameters.Login == "" || parameters.Password == "" {
			http.Error(w, "empty login/password", http.StatusBadRequest)
			return
		}

		hash, err := hash(parameters.Password)
		if err != nil {
			http.Error(w, "can't generate hash from password", http.StatusInternalServerError)
			return
		}

		conflict, err := s.db.Register(r.Context(), parameters.Login, hash)
		if err != nil {
			http.Error(w, "can't register", http.StatusInternalServerError)
			return
		}
		if conflict {
			http.Error(w, "login not unique", http.StatusConflict)
			return
		}

		token, err := authToken(parameters.Login)
		if err != nil {
			http.Error(w, "can't auth", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, token)

		w.WriteHeader(http.StatusOK)

	}
}

// Авторизация пользователя в приложении
func login(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var parameters model.RegisterRequest
		err := json.NewDecoder(r.Body).Decode(&parameters)
		if err != nil {
			http.Error(w, "Can't read body", http.StatusBadRequest)
			return
		}
		if parameters.Login == "" {
			http.Error(w, "empty login", http.StatusBadRequest)
			return
		}

		hash, err := s.db.GetPass(r.Context(), parameters.Login)
		if err != nil {
			http.Error(w, "can't auth", http.StatusInternalServerError)
			return
		}

		if !checkHash(parameters.Password, hash) {
			http.Error(w, "can't auth", http.StatusUnauthorized)
			return
		}

		token, err := authToken(parameters.Login)
		if err != nil {
			http.Error(w, "can't auth", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, token)

		w.WriteHeader(http.StatusOK)

	}
}
