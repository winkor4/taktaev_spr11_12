package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/winkor4/taktaev_spr11_12/internal/model"
)

func register(address string, creds model.Credentials) error {
	bodyData := registerRequest{
		Login:    creds.Login,
		Password: creds.Password,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(bodyData); err != nil {
		return err
	}
	resp, err := http.Post(address+"/user", "application/json", &buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("unexpected server status")
	}
	return nil
}

// func add(address string, creds model.Credentials, data model.Data) error {
// 	req, _ := http.NewRequest(http.MethodPost, address+"/api/content", &buf)
// 	req.SetBasicAuth(creds.Login, creds.Password)
// 	resp, err := http.DefaultClient.Do(req)
// 	// Обработка ответа
// 	return nil
// }
