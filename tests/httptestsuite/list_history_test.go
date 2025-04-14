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
	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/models"
)

func (s *APITestSuite) TestListHistory() {
	conf, err := config.New("config.yaml")
	s.Require().NoError(err)
	URL := fmt.Sprintf("http://%s/history/", net.JoinHostPort(conf.ConnectConfig.HTTP_HOST, conf.ConnectConfig.HTTP_PORT))
	s.Run("all good", func() {
		OrderID := int64(100)
		UserID := int64(1)
		Weight := int64(100)
		Price := int64(100)
		Pack := entity.Pack{PackType: "film", ExtraPack: "film"}
		Status := consts.ReturnedSt

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
		var orders []models.Order
		err = json.Unmarshal(body, &orders)
		s.Require().NoError(err)
		s.Assert().Equal(wantStatus, resp.StatusCode)
		s.Assert().Equal(UserID, orders[0].UserID)
		s.Assert().Equal(OrderID, orders[0].OrderID)
		s.Assert().Equal(Weight, orders[0].Weight)
		s.Assert().Equal(Price, orders[0].Price)
		s.Assert().Equal(Status, orders[0].Status)
		s.Assert().Equal(Pack.PackType, orders[0].Packaging.PackType)
		s.Assert().Equal(Pack.ExtraPack, orders[0].Packaging.ExtraPack)
	})
}
