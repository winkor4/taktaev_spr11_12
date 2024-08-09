package main

import (
	"fmt"
	"os"
)

const regCommand = "reg"

func main() {
	if len(os.Args) < 4 {
		fmt.Println("not enough arguments, expect login, password and command [options]")
		os.Exit(2)
	}
	// runAddress := os.Getenv("RUN_ADDRESS")

	// // Регистрируем команды
	// regCmd := flag.NewFlagSet(regCommand, flag.ContinueOnError)

	// // Зафиксированные позиции кредов: 1 и 2.
	// login, password := os.Args[1], os.Args[2]

	// // Зафиксированная позиция команды: 3
	// switch os.Args[3] {
	// case regCommand:
	// 	regCmd.Parse(os.Args[4:])
	// 	//internal/client/register({address, login, password})
	// 	//internal/client/cmd/Register(regCmd)
	// //...
	// default:
	// 	panic("unknown command")
	// }

}
