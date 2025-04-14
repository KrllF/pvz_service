//go:build integration

package httptest

import (
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

func TestReturnOrderToCourier(t *testing.T) {
	conf, err := config.New("config.yaml")
	require.NoError(t, err)
	t.Run("all good", func(t *testing.T) {
		URL := fmt.Sprintf("http://%s/return/100",
			net.JoinHostPort(conf.ConnectConfig.HTTP_HOST, conf.ConnectConfig.HTTP_PORT))
		tdb.SetUp(t, "Orders", "Users")
		defer tdb.TearDown(t)
		wantStatus := http.StatusOK
		wantBody := map[string]interface{}{
			"orderdel": 100,
		}

		wantBodyJSON, err := json.Marshal(wantBody)
		require.NoError(t, err)
		wantBodyJSON, err = json.MarshalIndent(wantBody, "", "\t")
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodDelete, URL, nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Basic "+basicAuth(conf.ConnectConfig.Login, conf.ConnectConfig.Password))
		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, wantStatus, resp.StatusCode)
		assert.Equal(t, string(wantBodyJSON), string(body))

		ctx := context.Background()

		query := `SELECT EXISTS (SELECT 1 FROM orders WHERE order_id = $1);`
		var exists bool
		require.NotEqual(t, tdb, nil)
		err = tdb.QueryRow(ctx, query, 100).Scan(&exists)
		require.NoError(t, err)
		assert.Equal(t, exists, false)
	})
	t.Run("bad request", func(t *testing.T) {
		URL := fmt.Sprintf("http://%s/return/2",
			net.JoinHostPort(conf.ConnectConfig.HTTP_HOST, conf.ConnectConfig.HTTP_PORT))
		tdb.SetUp(t, "Orders", "Users")
		defer tdb.TearDown(t)
		wantStatus := http.StatusNotFound
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodDelete, URL, nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Basic "+basicAuth(conf.ConnectConfig.Login, conf.ConnectConfig.Password))
		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, wantStatus, resp.StatusCode)
	})
}
