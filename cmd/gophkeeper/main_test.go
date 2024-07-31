// Тестироване основных функций приложения
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winkor4/taktaev_spr11_12/internal/crypto"
	"github.com/winkor4/taktaev_spr11_12/internal/log"
	"github.com/winkor4/taktaev_spr11_12/internal/pkg/config"
	"github.com/winkor4/taktaev_spr11_12/internal/server"
	"github.com/winkor4/taktaev_spr11_12/internal/storage"
)

// Параметры тестирования
type testParam struct {
	srv      *httptest.Server
	user     testUser
	masterSK string
}

// Тестовый пользователь
type testUser struct {
	login   string
	cookies []*http.Cookie
}

// Функция тестирования приложения
func TestApp(t *testing.T) {

	var parameters testParam
	parameters.srv = newTestSrv(t)
	parameters.masterSK = "[GUc7^q!u}!%RFGt"
	parameters.auth(t)
	parameters.addContentLogPass(t)
	parameters.getContentLogPass(t)
}

// Тестовый сервер
func newTestSrv(t *testing.T) *httptest.Server {

	ctx := context.Background()

	cfg, err := config.Parse()
	if err != nil {
		require.NoError(t, err)
	}

	db, err := storage.New(ctx, cfg.DatabaseURI)
	if err != nil {
		require.NoError(t, err)
	}

	logger, err := log.New()
	require.NoError(t, err)

	srv := server.New(server.Config{
		Cfg:    cfg,
		DB:     db,
		Logger: logger,
	})

	return httptest.NewServer(server.SrvRouter(srv))
}

// Регистрация/авторизация тестового пользователя
func (parameters *testParam) auth(t *testing.T) {

	type (
		reqSchema struct {
			Login    string `json:"login"`    // Логин
			Password string `json:"password"` // Пароль
		}
	)

	reqParam := reqSchema{
		Login:    "ivan",
		Password: "123",
	}

	byteParam, err := json.Marshal(reqParam)
	require.NoError(t, err)

	t.Run("Авторизация пользователя", func(t *testing.T) {

		body := bytes.NewReader(byteParam)
		request, err := http.NewRequest(http.MethodPost, parameters.srv.URL+"/auth", body)
		require.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")

		client := parameters.srv.Client()
		r, err := client.Do(request)
		require.NoError(t, err)

		if r.StatusCode == http.StatusUnauthorized {
			return
		}

		assert.Equal(t, http.StatusOK, r.StatusCode)

		parameters.user = testUser{
			login:   reqParam.Login,
			cookies: r.Cookies(),
		}

		err = r.Body.Close()
		require.NoError(t, err)

	})

	if parameters.user.login != "" {
		return
	}

	t.Run("Регистрация нового пользователя", func(t *testing.T) {

		body := bytes.NewReader(byteParam)
		request, err := http.NewRequest(http.MethodPost, parameters.srv.URL+"/user", body)
		require.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Key", parameters.masterSK)

		client := parameters.srv.Client()
		r, err := client.Do(request)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, r.StatusCode)

		parameters.user = testUser{
			login:   reqParam.Login,
			cookies: r.Cookies(),
		}

		err = r.Body.Close()
		require.NoError(t, err)

	})
}

// Загрузка данных типа логин/пароль
func (parameters *testParam) addContentLogPass(t *testing.T) {

	type (
		reqSchema struct {
			Name     string `json:"name"`     // Наименование
			Login    string `json:"login"`    // Логин
			Password string `json:"password"` // Пароль
		}
	)

	reqParam := reqSchema{
		Name:     "Моя почта",
		Login:    "mailLogin",
		Password: "mailPass",
	}

	byteParam, err := json.Marshal(reqParam)
	require.NoError(t, err)

	t.Run("Загрузка данных типа логин/пароль", func(t *testing.T) {

		body := bytes.NewReader(byteParam)
		request, err := http.NewRequest(http.MethodPost, parameters.srv.URL+"/api/content", body)
		require.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Data-Type", "LogPass")
		request.Header.Set("Key", parameters.masterSK)

		for _, c := range parameters.user.cookies {
			request.AddCookie(c)
		}

		client := parameters.srv.Client()
		r, err := client.Do(request)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, r.StatusCode)

		err = r.Body.Close()
		require.NoError(t, err)

	})
}

// Запрос данных типа логин/пароль
func (parameters *testParam) getContentLogPass(t *testing.T) {

	type (
		dataSchema struct {
			Name     string `json:"name"`     // Наименование
			Login    string `json:"login"`    // Логин
			Password string `json:"password"` // Пароль
		}
		response struct {
			Name   string `json:"name"`     // Наименование
			Data   string `json:"data"`     // Зашифрованные данные
			DataSK string `json:"data_key"` // Зашифрованный ключ данных
			EncSK  string `json:"key"`      // Зашифрованный ключ ключа данных
		}
	)

	wantData := dataSchema{
		Name:     "Моя почта",
		Login:    "mailLogin",
		Password: "mailPass",
	}

	byteWant, err := json.Marshal(wantData)
	require.NoError(t, err)

	t.Run("Запрос данных типа логин/пароль", func(t *testing.T) {

		name := "Моя почта"

		request, err := http.NewRequest(http.MethodGet, parameters.srv.URL+"/api/content/"+name, nil)
		require.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")

		for _, c := range parameters.user.cookies {
			request.AddCookie(c)
		}

		client := parameters.srv.Client()
		r, err := client.Do(request)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, r.StatusCode)

		rBody, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		err = r.Body.Close()
		require.NoError(t, err)

		var resp response
		err = json.Unmarshal(rBody, &resp)
		require.NoError(t, err)

		key, err := crypto.Decrypt(resp.EncSK, parameters.masterSK)
		require.NoError(t, err)
		dataKey, err := crypto.Decrypt(resp.DataSK, key)
		require.NoError(t, err)
		decData, err := crypto.Decrypt(resp.Data, dataKey)
		require.NoError(t, err)
		assert.JSONEq(t, string(byteWant), decData)

	})
}
