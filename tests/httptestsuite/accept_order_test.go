//go:build suite

package httptestsuite

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/KrllF/pvz_service/internal/config"
)

func (s *APITestSuite) TestAcceptOrder() {
	conf, err := config.New("config.yaml")
	s.Require().NoError(err)
	URL := fmt.Sprintf("http://%s/accept/", net.JoinHostPort(conf.ConnectConfig.HTTP_HOST, conf.ConnectConfig.HTTP_PORT))
	s.Run("all good", func() {
		s.T()
		order := map[string]interface{}{
			"order_id":   1001,
			"user_id":    100,
			"weight":     100,
			"price":      100,
			"pack_type":  "film",
			"extra_pack": "",
			"shelf_life": "2026-10-10T15:15:10Z",
		}
		wantStatus := http.StatusCreated
		wantBody := "Заказ успешно принят"
		requestBody, err := json.Marshal(order)
		s.Require().NoError(err)
		req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(requestBody))
		s.Require().NoError(err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Basic "+basicAuth(conf.ConnectConfig.Login, conf.ConnectConfig.Password))
		client := &http.Client{}
		resp, err := client.Do(req)
		s.Require().NoError(err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		s.Require().NoError(err)
		s.Assert().Equal(wantStatus, resp.StatusCode)
		s.Assert().Equal(wantBody, string(body))
		ctx := context.Background()

		query := `SELECT EXISTS (SELECT 1 FROM orders WHERE order_id = $1);`
		var exists bool
		s.Require().NotEqual(s, s.db, nil)
		err = s.db.QueryRow(ctx, query, 1001).Scan(&exists)
		s.Require().NoError(err)
		s.Assert().Equal(exists, true)
	})

	s.Run("bad request", func() {
		s.T()
		order := map[string]interface{}{
			"order_id":   1001,
			"user_id":    100,
			"weight":     100,
			"price":      100,
			"pack_type":  "film",
			"extra_pack": "fi",
			"shelf_life": "2026-10-10T15:15:10Z",
		}
		wantStatus := http.StatusNotFound

		requestBody, err := json.Marshal(order)
		s.Require().NoError(err)
		req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(requestBody))
		s.Require().NoError(err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Basic "+basicAuth(conf.ConnectConfig.Login, conf.ConnectConfig.Password))
		client := &http.Client{}
		resp, err := client.Do(req)
		s.Require().NoError(err)
		defer resp.Body.Close()
		s.Assert().Equal(wantStatus, resp.StatusCode)
	})
}
