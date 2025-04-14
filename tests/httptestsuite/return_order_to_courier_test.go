//go:build suite

package httptestsuite

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/KrllF/pvz_service/internal/config"
)

func (s *APITestSuite) TestReturnOrderToCourier() {
	conf, err := config.New("config.yaml")
	s.Require().NoError(err)
	s.Run("all good", func() {
		URL := fmt.Sprintf("http://%s/return/100",
			net.JoinHostPort(conf.ConnectConfig.HTTP_HOST, conf.ConnectConfig.HTTP_PORT))
		wantStatus := http.StatusOK
		wantBody := map[string]interface{}{
			"orderdel": 100,
		}

		wantBodyJSON, err := json.Marshal(wantBody)
		s.Require().NoError(err)
		wantBodyJSON, err = json.MarshalIndent(wantBody, "", "\t")
		s.Require().NoError(err)
		req, err := http.NewRequest(http.MethodDelete, URL, nil)
		s.Require().NoError(err)
		req.Header.Set("Authorization", "Basic "+basicAuth(conf.ConnectConfig.Login, conf.ConnectConfig.Password))
		client := &http.Client{}
		resp, err := client.Do(req)
		s.Require().NoError(err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		s.Require().NoError(err)
		s.Assert().Equal(wantStatus, resp.StatusCode)
		s.Assert().Equal(string(wantBodyJSON), string(body))

		ctx := context.Background()

		query := `SELECT EXISTS (SELECT 1 FROM orders WHERE order_id = $1);`
		var exists bool
		s.Require().NotEqual(s.pool, nil)
		err = s.pool.QueryRow(ctx, query, 100).Scan(&exists)
		s.Require().NoError(err)
		s.Assert().Equal(exists, false)
	})
	s.Run("bad request", func() {
		URL := fmt.Sprintf("http://%s/return/2",
			net.JoinHostPort(conf.ConnectConfig.HTTP_HOST, conf.ConnectConfig.HTTP_PORT))
		wantStatus := http.StatusNotFound
		s.Require().NoError(err)
		req, err := http.NewRequest(http.MethodDelete, URL, nil)
		s.Require().NoError(err)
		req.Header.Set("Authorization", "Basic "+basicAuth(conf.ConnectConfig.Login, conf.ConnectConfig.Password))
		client := &http.Client{}
		resp, err := client.Do(req)
		s.Require().NoError(err)
		defer resp.Body.Close()
		s.Assert().Equal(wantStatus, resp.StatusCode)
	})
}
