// Функции для запуска приложения
package app

import (
	"context"

	"github.com/winkor4/taktaev_spr11_12/internal/pkg/config"
	"github.com/winkor4/taktaev_spr11_12/internal/server"
	"github.com/winkor4/taktaev_spr11_12/internal/storage"
)

func Run() error {

	ctx := context.Background()

	cfg, err := config.Parse()
	if err != nil {
		return err
	}

	db, err := storage.New(ctx, cfg.DatabaseURI)
	if err != nil {
		return err
	}

	srv := server.New(server.Config{
		Cfg: cfg,
		DB:  db,
	})

	return srv.Run()

}
