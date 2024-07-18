// Описание типов данных
package model

// Входящие данные при регистрации
type RegisterRequest struct {
	Login    string `json:"login"`    // Логин
	Password string `json:"password"` // Пароль
}
