//go:build integration

package httptest

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/KrllF/pvz_service/internal/config"
	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/internal/handler/httph"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListOrdersUser(t *testing.T) {
	conf, err := config.New("config.yaml")
	require.NoError(t, err)
	URL := fmt.Sprintf("http://%s//history/orders/1?in_pvz=true", net.JoinHostPort(conf.ConnectConfig.HTTP_HOST, conf.ConnectConfig.HTTP_PORT))
	t.Run("all good", func(t *testing.T) {
		tdb.SetUp(t, "Orders", "Users")
		defer tdb.TearDown(t)
		UserID := int64(1)
		OrderID := int64(1)
		Weight := int64(100)
		Price := int64(100)
		Status := consts.DeliveredST

		wantStatus := http.StatusOK
		req, err := http.NewRequest(http.MethodGet, URL, nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Basic "+basicAuth(conf.ConnectConfig.Login, conf.ConnectConfig.Password))
		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		var orders httph.ListOrdersResponse
		err = json.Unmarshal(body, &orders)
		require.NoError(t, err)
		assert.Equal(t, wantStatus, resp.StatusCode)
		assert.Equal(t, UserID, orders.Orders[0].UserID)
		assert.Equal(t, OrderID, orders.Orders[0].OrderID)
		assert.Equal(t, Weight, orders.Orders[0].Weight)
		assert.Equal(t, Price, orders.Orders[0].Price)
		assert.Equal(t, Status, orders.Orders[0].Status)
	})
}
