package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/winkor4/taktaev_spr11_12/internal/content"
	"github.com/winkor4/taktaev_spr11_12/internal/users"
)

var errBadReq = errors.New("bad request")

// Регистрация пользователя в приложении
func addUser(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var userCredentials authSchema
		err := json.NewDecoder(r.Body).Decode(&userCredentials)
		if err != nil {
			http.Error(w, "Can't read body", http.StatusBadRequest)
			return
		}

		key := r.Header.Get("Key")
		if key == "" {
			http.Error(w, "can't read data type", http.StatusBadRequest)
			return
		}

		if userCredentials.Login == "" || userCredentials.Password == "" {
			http.Error(w, "empty login/password", http.StatusBadRequest)
			return
		}

		userManager := users.NewUserManager(userCredentials.Login,
			userCredentials.Password,
			key,
			s.db)
		conflict, err := userManager.AddUser(r.Context())
		switch {
		case err != nil:
			http.Error(w, "can't register", http.StatusInternalServerError)
			return
		case conflict:
			http.Error(w, "login not unique", http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// Авторизация пользователя в приложении
func atuhUser(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userCredentials authSchema
		err := json.NewDecoder(r.Body).Decode(&userCredentials)
		if err != nil {
			http.Error(w, "Can't read body", http.StatusBadRequest)
			return
		}
		if userCredentials.Login == "" {
			http.Error(w, "empty login", http.StatusBadRequest)
			return
		}

		auth, err := auth(r.Context(), userCredentials.Login,
			userCredentials.Password,
			"",
			s.db)
		switch {
		case err != nil:
			http.Error(w, "can't auth", http.StatusInternalServerError)
			return
		case !auth:
			http.Error(w, "can't auth", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// Запись данных на сервер
func addContent(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentManager, err := writreContentManager(r, s)
		if err != nil {
			http.Error(w, "can't read headers", http.StatusBadRequest)
			return
		}

		err = contentManager.AddContent(r.Context())
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

		var cfg content.Config
		cfg.User = user
		cfg.DB = s.db
		cfg.Name = name
		contentManager := content.NewContentManager(cfg)

		encData, err := contentManager.GetContent(r.Context())
		switch {
		case err != nil:
			http.Error(w, "can't get content by name", http.StatusInternalServerError)
			return
		case encData.ContentType == "":
			http.Error(w, "empty content", http.StatusNoContent)
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
		var cfg content.Config
		cfg.User = user
		cfg.DB = s.db
		contentManager := content.NewContentManager(cfg)

		dataList, err := contentManager.ContentList(r.Context())
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
		var cfg content.Config
		cfg.User = user
		cfg.DB = s.db
		cfg.Name = name
		contentManager := content.NewContentManager(cfg)

		err := contentManager.DeleteContent(r.Context())
		if err != nil {
			http.Error(w, "can't delete content from DB", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// Обновление данных на сервере
func updateContent(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentManager, err := writreContentManager(r, s)
		if err != nil {
			http.Error(w, "can't read headers", http.StatusBadRequest)
			return
		}
		err = contentManager.UpdateContent(r.Context())
		if err != nil {
			http.Error(w, "can't save content to DB", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	}
}

func writreContentManager(r *http.Request, s *Server) (*content.ContentManager, error) {
	contentManager := new(content.ContentManager)
	user, ok := userFromCtx(r.Context())
	if !ok {
		return contentManager, errBadReq
	}

	key := r.Header.Get("Key")
	dataType := r.Header.Get("Data-Type")
	if key == "" || dataType == "" {
		return contentManager, errBadReq
	}

	contentCredentials := getAddContentSchema(dataType)
	err := contentCredentials.JSONDecode(r.Body)
	if err != nil {
		return contentManager, err
	}

	var cfg content.Config
	cfg.User = user
	cfg.DB = s.db
	cfg.Key = key
	cfg.DataType = dataType
	cfg.ContentCredentials = contentCredentials
	contentManager = content.NewContentManager(cfg)
	return contentManager, nil
}
