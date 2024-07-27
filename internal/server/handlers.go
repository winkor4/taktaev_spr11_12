package server

import (
	"encoding/json"
	"net/http"
)

// Регистрация пользователя в приложении
func addUser(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var parameters authSchema
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

		conflict, err := s.db.AddUser(r.Context(), parameters.Login, hash)
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
			http.Error(w, "can't create Token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, token)

		w.WriteHeader(http.StatusOK)

	}
}

// Авторизация пользователя в приложении
func atuhUser(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var parameters authSchema
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

// Запись данных на сервер
func addContent(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// user, ok := userFromCtx(r.Context())
		// if !ok {
		// 	http.Error(w, "can't read login", http.StatusInternalServerError)
		// 	return
		// }

		dataType := r.Header.Get("Data-Type")
		if dataType == "" {
			http.Error(w, "can't read data type", http.StatusInternalServerError)
			return
		}

		parameters := getAddContentSchema(dataType)
		err := parameters.jsonDecode(r.Body)
		if err != nil {
			http.Error(w, "Can't read body", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)

	}
}
