//go:build suite

package httptestsuite

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/KrllF/pvz_service/internal/config"
	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/internal/handler/httph"
)

func (s *APITestSuite) TestListOrdersUser() {
	conf, err := config.New("config.yaml")
	s.Require().NoError(err)
	URL := fmt.Sprintf("http://%s//history/orders/1?in_pvz=true", net.JoinHostPort(conf.ConnectConfig.HTTP_HOST, conf.ConnectConfig.HTTP_PORT))
	s.Run("all good", func() {
		UserID := int64(1)
		OrderID := int64(1)
		Weight := int64(100)
		Price := int64(100)
		Status := consts.DeliveredST

		wantStatus := http.StatusOK
		req, err := http.NewRequest(http.MethodGet, URL, nil)
		s.Require().NoError(err)
		req.Header.Set("Authorization", "Basic "+basicAuth(conf.ConnectConfig.Login, conf.ConnectConfig.Password))
		client := &http.Client{}
		resp, err := client.Do(req)
		s.Require().NoError(err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		s.Require().NoError(err)
		var orders httph.ListOrdersResponse
		err = json.Unmarshal(body, &orders)
		s.Require().NoError(err)
		s.Assert().Equal(wantStatus, resp.StatusCode)
		s.Assert().Equal(UserID, orders.Orders[0].UserID)
		s.Assert().Equal(OrderID, orders.Orders[0].OrderID)
		s.Assert().Equal(Weight, orders.Orders[0].Weight)
		s.Assert().Equal(Price, orders.Orders[0].Price)
		s.Assert().Equal(Status, orders.Orders[0].Status)
	})
}
