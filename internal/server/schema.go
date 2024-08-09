package server

import (
	"encoding/json"
	"io"

	"github.com/google/uuid"
	"github.com/winkor4/taktaev_spr11_12/internal/crypto"
	"github.com/winkor4/taktaev_spr11_12/internal/model"
)

// Входящие данные при регистрации и авторизации
type authSchema struct {
	Login    string `json:"login"`    // Логин
	Password string `json:"password"` // Пароль
}

// Данные возвращаемые сервером
type getContentResponse struct {
	Name     string `json:"name"`      // Наименование
	DataType string `json:"data_type"` // Тип данных
	Data     string `json:"data"`      // Зашифрованные данные
	DataSK   string `json:"data_key"`  // Зашифрованный ключ данных
	EncSK    string `json:"key"`       // Зашифрованный ключ ключа данных
}

// Данные возвращаемые сервером
type contentListResponse struct {
	Name string `json:"name"` // Наименование
}

// Описание данных типа логин/пароль
type addContentLogPass struct {
	Name     string `json:"name"`     // Наименование
	Login    string `json:"login"`    // Логин
	Password string `json:"password"` // Пароль
}

func getAddContentSchema(dataType string) model.AddContentRequest {

	var result model.AddContentRequest

	if dataType == "LogPass" {
		result = &addContentLogPass{}
	}

	return result
}

func (schema *addContentLogPass) JSONDecode(body io.ReadCloser) error {
	err := json.NewDecoder(body).Decode(schema)
	if err != nil {
		return err
	}
	return nil
}

func (schema *addContentLogPass) SchemaToStorageData(user model.User, dataType string) (model.StorageData, error) {

	var result model.StorageData

	data, err := json.Marshal(schema)
	if err != nil {
		return result, err
	}

	dataSK := crypto.RandStr(16)
	encData, err := crypto.Encrypt(string(data), dataSK)
	if err != nil {
		return result, err
	}
	encKey, err := crypto.Encrypt(dataSK, user.Key)
	if err != nil {
		return result, err
	}

	result.ID = uuid.New().String()
	result.User = user
	result.Name = schema.Name
	result.ContentType = dataType
	result.Data = encData
	result.DataSK = encKey

	return result, nil
}

func encDataToSchema(data model.EncContent) getContentResponse {
	return getContentResponse{
		Name:     data.Name,
		DataType: data.ContentType,
		Data:     data.Data,
		DataSK:   data.DataSK,
		EncSK:    data.EncSK,
	}
}

func dataListToSchema(dataList []string) []contentListResponse {
	result := make([]contentListResponse, 0, len(dataList))
	for _, v := range dataList {
		var resp contentListResponse
		resp.Name = v
		result = append(result, resp)
	}
	return result
}
