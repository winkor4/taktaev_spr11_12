package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/winkor4/taktaev_spr11_12/internal/crypto"
	"github.com/winkor4/taktaev_spr11_12/internal/model"
)

// Регистрация
func (c *Client) register() error {
	bodyData := registerRequest{
		Login:    c.login,
		Password: c.password,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(bodyData); err != nil {
		return err
	}
	key := crypto.PasswordHash(c.password)

	req, err := http.NewRequest(http.MethodPost, c.runAddress+"/user", &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Key", key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println("Успешная регистрация")
		return nil
	case http.StatusConflict:
		return errors.New("указанный логин уже занят")
	case http.StatusBadRequest:
		return errors.New("неверные параметры запроса")
	default:
		return errors.New("не удалось зарегестрироваться, ошибка на сервере")
	}
}

// Проверка авторизации
func (c *Client) auth() error {
	bodyData := registerRequest{
		Login:    c.login,
		Password: c.password,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(bodyData); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.runAddress+"/auth", &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := new(http.Client)
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println("Проверка авторизации пройдена")
		return nil
	case http.StatusUnauthorized:
		return errors.New("неверный логин или пароль")
	case http.StatusBadRequest:
		return errors.New("неверные параметры запроса")
	default:
		return errors.New("не удалось пройти авторизацию, ошибка на сервере")
	}
}

// Отправляет контент на сервер
func (c *Client) addContent(reqData model.AddRequest) error {
	body := bytes.NewReader(reqData.Body)
	req, err := http.NewRequest(http.MethodPost, c.runAddress+"/api/content", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Data-Type", reqData.DataType)
	req.Header.Set("Key", crypto.PasswordHash(c.password))
	req.SetBasicAuth(c.login, c.password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println("Данные успешно сохранены")
		return nil
	case http.StatusUnauthorized:
		return errors.New("неверный логин или пароль")
	case http.StatusBadRequest:
		return errors.New("неверные параметры запроса")
	default:
		return errors.New("не удалось сохранить данные, ошибка на сервере")
	}
}

// Получает контент от сервера по имени
func (c *Client) getContent(name string) error {
	req, err := http.NewRequest(http.MethodGet, c.runAddress+"/api/content/"+name, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.login, c.password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Printf("Данные по запросу %s:\n", name)
	case http.StatusNoContent:
		fmt.Println("По заданному имени ничего не найдено")
		return nil
	case http.StatusUnauthorized:
		return errors.New("неверный логин или пароль")
	default:
		return errors.New("не удалось найти данные, ошибка на сервере")
	}

	rBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var content getContentResponse
	err = json.Unmarshal(rBody, &content)
	if err != nil {
		return err
	}

	masterSK := crypto.PasswordHash(c.password)
	key, err := crypto.Decrypt(content.EncSK, masterSK)
	if err != nil {
		return err
	}
	dataKey, err := crypto.Decrypt(content.DataSK, key)
	if err != nil {
		return err
	}
	decData, err := crypto.Decrypt(content.Data, dataKey)
	if err != nil {
		return err
	}
	fmt.Println(decData)
	return nil
}

// Удаляет контент на сервере
func (c *Client) delContent(name string) error {
	req, err := http.NewRequest(http.MethodDelete, c.runAddress+"/api/content/"+name, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.login, c.password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println("Данные успешно удалены")
		return nil
	case http.StatusUnauthorized:
		return errors.New("неверный логин или пароль")
	default:
		return errors.New("не удалось найти данные, ошибка на сервере")
	}
}

// Обновляет контент на сервере
func (c *Client) updateContent(reqData model.AddRequest) error {
	body := bytes.NewReader(reqData.Body)
	req, err := http.NewRequest(http.MethodPost, c.runAddress+"/api/content/update", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Data-Type", reqData.DataType)
	req.Header.Set("Key", crypto.PasswordHash(c.password))
	req.SetBasicAuth(c.login, c.password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println("Данные успешно обновлены")
		return nil
	case http.StatusUnauthorized:
		return errors.New("неверный логин или пароль")
	case http.StatusBadRequest:
		return errors.New("неверные параметры запроса")
	default:
		return errors.New("не удалось сохранить данные, ошибка на сервере")
	}
}

// Получает список контента от сервера
func (c *Client) listContent() error {
	req, err := http.NewRequest(http.MethodGet, c.runAddress+"/api/content", nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.login, c.password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Printf("Данные по запросу:\n")
	case http.StatusUnauthorized:
		return errors.New("неверный логин или пароль")
	default:
		return errors.New("не удалось найти данные, ошибка на сервере")
	}

	rBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	content := make([]listContentResponse, 0)
	err = json.Unmarshal(rBody, &content)
	if err != nil {
		return err
	}

	for _, v := range content {
		fmt.Printf("%s\n", v.Name)
	}
	return nil
}
