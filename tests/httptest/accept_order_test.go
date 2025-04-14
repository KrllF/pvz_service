//go:build integration

package httptest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/KrllF/pvz_service/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAcceptOrder(t *testing.T) {
	conf, err := config.New("config.yaml")
	require.NoError(t, err)
	URL := fmt.Sprintf("http://%s/accept/", net.JoinHostPort(conf.ConnectConfig.HTTP_HOST, conf.ConnectConfig.HTTP_PORT))
	t.Run("all good", func(t *testing.T) {
		tdb.SetUp(t, "Orders", "Users")
		defer tdb.TearDown(t)
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
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Basic "+basicAuth(conf.ConnectConfig.Login, conf.ConnectConfig.Password))
		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, wantStatus, resp.StatusCode)
		assert.Equal(t, wantBody, string(body))

		ctx := context.Background()

		query := `SELECT EXISTS (SELECT 1 FROM orders WHERE order_id = $1);`
		var exists bool
		require.NotEqual(t, tdb, nil)
		err = tdb.QueryRow(ctx, query, 1001).Scan(&exists)
		require.NoError(t, err)
		assert.Equal(t, exists, true)
	})
	t.Run("bad request", func(t *testing.T) {
		tdb.SetUp(t, "Orders", "Users")
		defer tdb.TearDown(t)
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
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Basic "+basicAuth(conf.ConnectConfig.Login, conf.ConnectConfig.Password))
		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, wantStatus, resp.StatusCode)
	})
}
