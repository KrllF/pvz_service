//go:build suite

package httptestsuite

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/KrllF/pvz_service/internal/config"
	"github.com/KrllF/pvz_service/internal/handler/httph"
)

func (s *APITestSuite) TestProcessClient() {
	conf, err := config.New("config.yaml")
	s.Require().NoError(err)
	URL := fmt.Sprintf("http://%s/process", net.JoinHostPort(conf.ConnectConfig.HTTP_HOST, conf.ConnectConfig.HTTP_PORT))
	s.Run("all good", func() {
		reqBody := httph.ProcessOrderRequest{
			UserID:    1,
			Operation: "issue",
			OrderIDs:  []int64{1},
		}
		wantStatus := http.StatusOK
		requestBody, err := json.Marshal(reqBody)
		s.Require().NoError(err)
		req, err := http.NewRequest(http.MethodPut, URL, bytes.NewBuffer(requestBody))
		s.Require().NoError(err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Basic "+basicAuth(conf.ConnectConfig.Login, conf.ConnectConfig.Password))
		client := &http.Client{}
		resp, err := client.Do(req)
		s.Require().NoError(err)
		defer resp.Body.Close()
		s.Require().NoError(err)
		s.Assert().Equal(wantStatus, resp.StatusCode)

		ctx := context.Background()

		query := `SELECT status_id FROM orders WHERE order_id=$1;`
		var exists int64
		s.Require().NoError(err)
		err = s.pool.QueryRow(ctx, query, 1).Scan(&exists)
		s.Require().NoError(err)
		s.Assert().Equal(exists, int64(2))
	})
}
