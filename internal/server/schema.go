package server

import (
	"encoding/json"
	"io"

	"github.com/winkor4/taktaev_spr11_12/internal/model"
)

// Входящие данные при регистрации и авторизации
type authSchema struct {
	Login    string `json:"login"`    // Логин
	Password string `json:"password"` // Пароль
}

// Интерфейс, который ползволяет прочитать тело запроса в зависимости от типа данных
type addContentRequest interface {
	jsonDecode(body io.ReadCloser) error
}

// Описание данных типа логин/пароль
type addContentLogPass struct {
	Name     string `json:"name"`     // Наименование
	Login    string `json:"login"`    // Логин
	Password string `json:"password"` // Пароль
}

func getAddContentSchema(dataType string) addContentRequest {

	var result addContentRequest

	if dataType == "LogPass" {
		result = &addContentLogPass{}
	}

	return result

}

func (schema *addContentLogPass) jsonDecode(body io.ReadCloser) error {
	err := json.NewDecoder(body).Decode(schema)
	if err != nil {
		return err
	}
	return nil
}

func addUserReqToModel(l string, p string, k string) model.User {
	return model.User{
		Login:    l,
		Password: p,
		Key:      k,
	}
}
