package server

import (
	"encoding/json"
	"io"
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
