// Тестироване основных функций приложения
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	login    string
	password string
}

// Функция тестирования приложения
func TestApp(t *testing.T) {

	var parameters testParam
	parameters.srv = newTestSrv(t)
	parameters.masterSK = "[GUc7^q!u}!%RFGt"
	parameters.auth(t)
	parameters.addContentLogPass(t)
	parameters.updateContent(t)
	parameters.getContentLogPass(t)
	parameters.deleteContent(t)
	parameters.сontentList(t)
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
		Password: "asduq6we79gq7wfd",
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
			login:    reqParam.Login,
			password: reqParam.Password,
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
			login:    reqParam.Login,
			password: reqParam.Password,
		}

		err = r.Body.Close()
		require.NoError(t, err)

	})
}

// Загрузка данных типа логин/пароль
func (parameters *testParam) addContentLogPass(t *testing.T) {

	type (
		testAPP struct {
			testName  string
			byteParam []byte
		}
		reqSchema struct {
			Name     string `json:"name"`     // Наименование
			Login    string `json:"login"`    // Логин
			Password string `json:"password"` // Пароль
		}
	)

	reqSlice := []reqSchema{
		{
			Name:     "Моя почта",
			Login:    "mailLogin",
			Password: "mailPass",
		},
		{
			Name:     "Моя рабочая почта",
			Login:    "mailLogin",
			Password: "mailPass",
		},
		{
			Name:     "Почта для обновления",
			Login:    "mailLogin",
			Password: "mailPass",
		},
	}

	byteSlice := make([][]byte, 3)

	for i, v := range reqSlice {
		byteParam, err := json.Marshal(v)
		require.NoError(t, err)
		byteSlice[i] = byteParam
	}

	testTable := []testAPP{
		{
			testName:  "Загрузка данных типа логин/пароль",
			byteParam: byteSlice[0],
		},
		{
			testName:  "Загрузка данных типа логин/пароль для удаления",
			byteParam: byteSlice[1],
		},
		{
			testName:  "Загрузка данных типа логин/пароль для обновления",
			byteParam: byteSlice[2],
		},
	}

	for _, testData := range testTable {
		t.Run(testData.testName, func(t *testing.T) {

			body := bytes.NewReader(testData.byteParam)
			request, err := http.NewRequest(http.MethodPost, parameters.srv.URL+"/api/content", body)
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Data-Type", "LogPass")
			request.Header.Set("Key", parameters.masterSK)
			request.SetBasicAuth(parameters.user.login, parameters.user.password)

			client := parameters.srv.Client()
			r, err := client.Do(request)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, r.StatusCode)

			err = r.Body.Close()
			require.NoError(t, err)

		})
	}
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
			Name     string `json:"name"`      // Наименование
			DataType string `json:"data_type"` // Тип данных
			Data     string `json:"data"`      // Зашифрованные данные
			DataSK   string `json:"data_key"`  // Зашифрованный ключ данных
			EncSK    string `json:"key"`       // Зашифрованный ключ ключа данных
		}
	)

	wantData := dataSchema{
		Name:     "Почта для обновления",
		Login:    "NEWmailLogin",
		Password: "NEWmailPass",
	}

	byteWant, err := json.Marshal(wantData)
	require.NoError(t, err)

	t.Run("Запрос данных типа логин/пароль", func(t *testing.T) {

		request, err := http.NewRequest(http.MethodGet, parameters.srv.URL+"/api/content/"+wantData.Name, nil)
		require.NoError(t, err)

		request.SetBasicAuth(parameters.user.login, parameters.user.password)

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
		assert.Equal(t, "LogPass", resp.DataType)

	})
}

// Удаление ранее загруженных данных
func (parameters *testParam) deleteContent(t *testing.T) {
	t.Run("Запрос данных типа логин/пароль", func(t *testing.T) {

		name := "Моя рабочая почта"

		request, err := http.NewRequest(http.MethodDelete, parameters.srv.URL+"/api/content/"+name, nil)
		require.NoError(t, err)

		request.SetBasicAuth(parameters.user.login, parameters.user.password)

		client := parameters.srv.Client()
		r, err := client.Do(request)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, r.StatusCode)

		err = r.Body.Close()
		require.NoError(t, err)

	})
}

// Обновление данных типа логин/пароль
func (parameters *testParam) updateContent(t *testing.T) {

	type (
		testAPP struct {
			testName  string
			byteParam []byte
		}
		reqSchema struct {
			Name     string `json:"name"`     // Наименование
			Login    string `json:"login"`    // Логин
			Password string `json:"password"` // Пароль
		}
	)

	reqSlice := []reqSchema{
		{
			Name:     "Почта для обновления",
			Login:    "NEWmailLogin",
			Password: "NEWmailPass",
		},
	}

	byteSlice := make([][]byte, 1)

	for i, v := range reqSlice {
		byteParam, err := json.Marshal(v)
		require.NoError(t, err)
		byteSlice[i] = byteParam
	}

	testTable := []testAPP{
		{
			testName:  "Обновление данных на сервере",
			byteParam: byteSlice[0],
		},
	}

	for _, testData := range testTable {
		t.Run(testData.testName, func(t *testing.T) {

			body := bytes.NewReader(testData.byteParam)
			request, err := http.NewRequest(http.MethodPost, parameters.srv.URL+"/api/content/update", body)
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Data-Type", "LogPass")
			request.Header.Set("Key", parameters.masterSK)

			request.SetBasicAuth(parameters.user.login, parameters.user.password)

			client := parameters.srv.Client()
			r, err := client.Do(request)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, r.StatusCode)

			err = r.Body.Close()
			require.NoError(t, err)

		})
	}
}

// Запрос списка загруженных данных
func (parameters *testParam) сontentList(t *testing.T) {

	// Данные возвращаемые сервером
	type response struct {
		Name string `json:"name"` // Наименование
	}

	want := make(map[string]bool, 0)
	want["Моя почта"] = true
	want["Почта для обновления"] = true

	t.Run("Запрос списка загруженных данных", func(t *testing.T) {

		request, err := http.NewRequest(http.MethodGet, parameters.srv.URL+"/api/content", nil)
		require.NoError(t, err)

		request.SetBasicAuth(parameters.user.login, parameters.user.password)

		client := parameters.srv.Client()
		r, err := client.Do(request)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, r.StatusCode)

		rBody, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		err = r.Body.Close()
		require.NoError(t, err)

		resp := make([]response, 0)
		err = json.Unmarshal(rBody, &resp)
		require.NoError(t, err)

		for _, v := range resp {
			_, ok := want[v.Name]
			if !ok {
				require.NoError(t, errors.New("no content"))
			}
		}

	})
}
