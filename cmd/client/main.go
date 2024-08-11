package main

import (
	"fmt"
	"os"

	"github.com/winkor4/taktaev_spr11_12/internal/client"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Не достаточно параметров, ожидаются параметры login, password, command")
		os.Exit(2)
	}

	// Зафиксированные позиции кредов: 1 и 2.
	login, password := os.Args[1], os.Args[2]
	if login == "" || password == "" {
		fmt.Println("Логин и пароль не могут быть пустыми")
		os.Exit(2)
	}

	// Зафиксированная позиция команды: 3
	command := os.Args[3]
	if command == "" {
		fmt.Println("Имя команды не может быть пустой")
		os.Exit(2)
	}

	cfg := client.Config{
		Login:      login,
		Password:   password,
		RunAddress: "http://localhost:8080",
	}

	// cfg := client.Config{
	// 	Login:      "new",
	// 	Password:   "123",
	// 	RunAddress: "http://localhost:8080",
	// }

	client := client.NewClient(cfg)
	err := client.Do(command)
	if err != nil {
		fmt.Println("Ошибка выполнения команды:")
		fmt.Println(err)
		os.Exit(2)
	}

}
