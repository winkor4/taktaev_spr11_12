// Тестироване основных функций приложения
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	parameters.masterSK = "abc&1*~#^2^#s0^=)^^"
	parameters.auth(t)
	parameters.addContentLogPass(t)
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
