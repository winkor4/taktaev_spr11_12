// Описание типов данных
package model

// Входящие данные при регистрации
type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
