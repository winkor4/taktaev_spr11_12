package client

import (
	"encoding/json"
	"errors"
	"flag"
	"strings"

	"github.com/winkor4/taktaev_spr11_12/internal/model"
)

func readAddFlags() (model.AddRequest, error) {
	var (
		result   model.AddRequest
		dataType string
	)

	flag.Parse()
	addCmd := flag.NewFlagSet("add", flag.ContinueOnError)
	addCmd.StringVar(&dataType, "type", "", "тип данных")
	if err := addCmd.Parse(flag.Args()[3:4]); err != nil {
		return result, err
	}
	if dataType == "" {
		return result, errors.New("тип данных не может быть пустым")
	}
	dataType = strings.ToLower(dataType)
	included, ok := dataTypes[dataType]
	switch {
	case !ok:
		return result, errors.New("неизвестный тип данных")
	case !included:
		return result, errors.New("данный тип данных больше не поддерживается")
	}

	var err error
	switch dataType {
	case "logpass":
		result.DataType = "LogPass"
		err = readLogPassFlags(&result, addCmd)
	}
	if err != nil {
		return result, err
	}

	return result, nil
}

func readLogPassFlags(result *model.AddRequest, sub *flag.FlagSet) error {
	var (
		name     string
		login    string
		password string
	)
	sub.StringVar(&name, "name", "", "имя сервиса")
	sub.StringVar(&login, "l", "", "логин")
	sub.StringVar(&password, "p", "", "пароль")
	if err := sub.Parse(flag.Args()[4:]); err != nil {
		return err
	}
	if name == "" || login == "" || password == "" {
		return errors.New("должны быть заполнены параметры name, login и password")
	}
	schema := addLogPassRequest{
		Name:     name,
		Login:    login,
		Password: password,
	}
	body, err := json.Marshal(schema)
	if err != nil {
		return err
	}
	result.Body = body
	return nil
}

func readGetDelFlags() (string, error) {
	var name string
	flag.Parse()
	addCmd := flag.NewFlagSet("get", flag.ContinueOnError)
	addCmd.StringVar(&name, "name", "", "имя сервиса")
	if err := addCmd.Parse(flag.Args()[3:]); err != nil {
		return "", err
	}
	if name == "" {
		return "", errors.New("имя сервиса не может быть пустым")
	}
	return name, nil
}
