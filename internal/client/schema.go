package client

type registerRequest struct {
	Login    string `json:"login"`    // Логин
	Password string `json:"password"` // Пароль
}

// Описание данных типа логин/пароль
type addLogPassRequest struct {
	Name     string `json:"name"`     // Наименование
	Login    string `json:"login"`    // Логин
	Password string `json:"password"` // Пароль
}

// Данные возвращаемые сервером getContent
type getContentResponse struct {
	Name     string `json:"name"`      // Наименование
	DataType string `json:"data_type"` // Тип данных
	Data     string `json:"data"`      // Зашифрованные данные
	DataSK   string `json:"data_key"`  // Зашифрованный ключ данных
	EncSK    string `json:"key"`       // Зашифрованный ключ ключа данных
}

// Данные возвращаемые сервером listContent
type listContentResponse struct {
	Name string `json:"name"` // Наименование
}
