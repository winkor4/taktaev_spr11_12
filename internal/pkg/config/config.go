// Модуль config парсит конфигурацию для сервера.
package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

// Config параметры конфигурации
type Config struct {
	RunAddress  string `env:"RUN_ADDRESS"`
	DatabaseURI string `env:"DATABASE_GOPHKEEPER"`
}

var (
	flagRunAddress  string // Адрес сервера
	flagDatabaseURI string // Адрес базы данных
)

// Parse - возвращает конфигурацию для сервера
func Parse() (*Config, error) {

	cfg := new(Config)

	flag.StringVar(&flagRunAddress, "a", "", "srv run addres and port")
	flag.StringVar(&flagDatabaseURI, "d", "", "PostgresSQL server")
	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	if cfg.RunAddress == "" {
		cfg.RunAddress = flagRunAddress
	}
	if cfg.DatabaseURI == "" {
		cfg.DatabaseURI = flagDatabaseURI
	}

	return cfg, nil

}
