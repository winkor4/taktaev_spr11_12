// Функции для запуска и работы сервера
package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/winkor4/taktaev_spr11_12/internal/pkg/config"
	"github.com/winkor4/taktaev_spr11_12/internal/storage"
)

// Config - параметры создания сервера
type Config struct {
	Cfg *config.Config
	DB  *storage.DB
}

// Server - описание сервера
type Server struct {
	cfg *config.Config
	db  *storage.DB
}

// New - возвращает новый сервер
func New(cfg Config) *Server {
	return &Server{
		cfg: cfg.Cfg,
		db:  cfg.DB,
	}
}

// Run - запускает сервер
func (s *Server) Run() error {
	return http.ListenAndServe(s.cfg.RunAddress, SrvRouter(s))
}

// SrvRouter - возвращает новый объект Mux
func SrvRouter(s *Server) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/api/register", checkContentType(register(s), "application/json"))
	r.Post("/api/login", checkContentType(login(s), "application/json"))
	r.Mount("/api/user", userRouter(s))

	return r
}

// UserRouter - возвращает новый объект Mux
func userRouter(s *Server) *chi.Mux {
	r := chi.NewRouter()
	r.Use(authorization())

	r.Post("/data/text", checkContentType(uploadTextData(s), "application/json"))

	return r
}
