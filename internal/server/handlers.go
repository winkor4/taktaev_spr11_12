package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/winkor4/taktaev_spr11_12/internal/crypto"
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

		key := r.Header.Get("Key")
		if key == "" {
			http.Error(w, "can't read data type", http.StatusBadRequest)
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

		encryptionSK := crypto.RandStr(16)
		encryptionSK, err = crypto.Encrypt(encryptionSK, key)
		if err != nil {
			http.Error(w, "can't Encrypt key", http.StatusInternalServerError)
			return
		}

		conflict, err := s.db.AddUser(r.Context(), gerUserModel(parameters.Login, hash, encryptionSK))
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

		user, ok := userFromCtx(r.Context())
		if !ok {
			http.Error(w, "can't read login", http.StatusInternalServerError)
			return
		}

		key := r.Header.Get("Key")
		if key == "" {
			http.Error(w, "can't read data type", http.StatusBadRequest)
			return
		}

		dataType := r.Header.Get("Data-Type")
		if dataType == "" {
			http.Error(w, "can't read data type", http.StatusBadRequest)
			return
		}

		parameters := getAddContentSchema(dataType)
		err := parameters.jsonDecode(r.Body)
		if err != nil {
			http.Error(w, "Can't read body", http.StatusBadRequest)
			return
		}

		encKey, err := s.db.GetKey(r.Context(), user)
		if err != nil {
			http.Error(w, "can't auth", http.StatusInternalServerError)
			return
		}

		key, err = crypto.Decrypt(encKey, key)
		if err != nil {
			http.Error(w, "can't Decrypt key", http.StatusInternalServerError)
			return
		}

		sData, err := parameters.schemaToStorageData(gerUserModel(user, "", key), dataType)
		if err != nil {
			http.Error(w, "can't save data", http.StatusInternalServerError)
			return
		}

		err = s.db.AddContent(r.Context(), sData)
		if err != nil {
			http.Error(w, "can't save content to DB", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	}
}

// Возвращает данные по имени
func getContent(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		user, ok := userFromCtx(r.Context())
		if !ok {
			http.Error(w, "can't read login", http.StatusInternalServerError)
			return
		}
		encData, err := s.db.GetContent(r.Context(), name, user)
		if err != nil {
			http.Error(w, "can't get content from DB", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(encDataToSchema(encData)); err != nil {
			http.Error(w, "Can't encode response", http.StatusInternalServerError)
			return
		}
	}
}

// Возвращает список данных
func contentList(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := userFromCtx(r.Context())
		if !ok {
			http.Error(w, "can't read login", http.StatusInternalServerError)
			return
		}
		dataList, err := s.db.ContentList(r.Context(), user)
		if err != nil {
			http.Error(w, "can't get content from DB", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(dataListToSchema(dataList)); err != nil {
			http.Error(w, "Can't encode response", http.StatusInternalServerError)
			return
		}
	}
}

// Удаляет данные с сервера
func deleteContent(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		user, ok := userFromCtx(r.Context())
		if !ok {
			http.Error(w, "can't read login", http.StatusInternalServerError)
			return
		}
		err := s.db.DeleteContent(r.Context(), name, user)
		if err != nil {
			http.Error(w, "can't delete content from DB", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
