// Функции для запуска и работы сервера
package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/winkor4/taktaev_spr11_12/internal/log"
	"github.com/winkor4/taktaev_spr11_12/internal/pkg/config"
	"github.com/winkor4/taktaev_spr11_12/internal/storage"
	"go.uber.org/zap/zapcore"
)

// Config - параметры создания сервера
type Config struct {
	Cfg    *config.Config
	DB     *storage.DB
	Logger *log.Logger
}

// Server - описание сервера
type Server struct {
	cfg    *config.Config
	db     *storage.DB
	logger *log.Logger
}

// New - возвращает новый сервер
func New(cfg Config) *Server {
	return &Server{
		cfg:    cfg.Cfg,
		db:     cfg.DB,
		logger: cfg.Logger,
	}
}

// Run - запускает сервер
func (s *Server) Run() error {
	s.logger.Logw(zapcore.InfoLevel, "Starting server", "Address", s.cfg.RunAddress)
	return http.ListenAndServe(s.cfg.RunAddress, SrvRouter(s))
}

// SrvRouter - возвращает новый объект Mux
func SrvRouter(s *Server) *chi.Mux {
	r := chi.NewRouter()
	r.Use(logMiddleware(s))

	r.Post("/user", checkContentType(addUser(s), "application/json"))
	r.Post("/auth", checkContentType(atuhUser(s), "application/json"))
	r.Mount("/api", apiRouter(s))

	return r
}

// apiRouter - возвращает новый объект Mux
func apiRouter(s *Server) *chi.Mux {
	r := chi.NewRouter()
	r.Use(authorization())

	r.Post("/content", checkContentType(addContent(s), "application/json"))

	return r
}
