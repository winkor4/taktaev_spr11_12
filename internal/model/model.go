// Описание типов данных
package model

// Входящие данные при регистрации
type RegisterRequest struct {
	Login    string `json:"login"`    // Логин
	Password string `json:"password"` // Пароль
}

// Список текстовых данных для записи в БД
type TextDataList struct {
	ID   string `json:"id"`   // Идентификатор данных
	Data string `json:"data"` // Произвольные текстовые данные
}

type UploadTextDataResponse struct {
	ID       string `json:"id"`       // Идентификатор данных
	Conflict bool   `json:"conflict"` // Идентификатор уже занят
}
