// Описание типов данных
package model

// Описание пользователя
type User struct {
	Login    string // Логин
	Password string // Пароль
	Key      string // Ключ шифрования
}

// Описание данных
type StorageData struct {
	ID     string // UUID
	Name   string // Наименование
	User   User   // Пользователь
	Data   string // Зашифрованные данные
	DataSK string // Зашифрованный ключ
}
