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
	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListHistory(t *testing.T) {
	conf, err := config.New("config.yaml")
	require.NoError(t, err)
	URL := fmt.Sprintf("http://%s/history/", net.JoinHostPort(conf.ConnectConfig.HTTP_HOST, conf.ConnectConfig.HTTP_PORT))
	t.Run("all good", func(t *testing.T) {
		tdb.SetUp(t, "Orders", "Users")
		defer tdb.TearDown(t)
		OrderID := int64(100)
		UserID := int64(1)
		Weight := int64(100)
		Price := int64(100)
		Pack := entity.Pack{PackType: "film", ExtraPack: "film"}
		Status := consts.ReturnedSt

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
		var orders []models.Order
		err = json.Unmarshal(body, &orders)
		require.NoError(t, err)
		assert.Equal(t, wantStatus, resp.StatusCode)
		assert.Equal(t, UserID, orders[0].UserID)
		assert.Equal(t, OrderID, orders[0].OrderID)
		assert.Equal(t, Weight, orders[0].Weight)
		assert.Equal(t, Price, orders[0].Price)
		assert.Equal(t, Status, orders[0].Status)
		assert.Equal(t, Pack.PackType, orders[0].Packaging.PackType)
		assert.Equal(t, Pack.ExtraPack, orders[0].Packaging.ExtraPack)
	})
}
