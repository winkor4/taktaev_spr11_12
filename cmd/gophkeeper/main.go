// Модуль запуска приложения.
package main

import (
	"log"

	"github.com/winkor4/taktaev_spr11_12/internal/pkg/app"
)

// main - старт приложения
func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
