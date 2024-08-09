package client

type registerRequest struct {
	Login    string `json:"login"`    // Логин
	Password string `json:"password"` // Пароль
}
