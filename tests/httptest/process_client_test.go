//go:build integration

package httptest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/KrllF/pvz_service/internal/config"
	"github.com/KrllF/pvz_service/internal/handler/httph"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessClient(t *testing.T) {
	conf, err := config.New("config.yaml")
	require.NoError(t, err)
	URL := fmt.Sprintf("http://%s/process", net.JoinHostPort(conf.ConnectConfig.HTTP_HOST, conf.ConnectConfig.HTTP_PORT))
	t.Run("all good", func(t *testing.T) {
		tdb.SetUp(t, "Orders", "Users")
		defer tdb.TearDown(t)
		reqBody := httph.ProcessOrderRequest{
			UserID:    1,
			Operation: "issue",
			OrderIDs:  []int64{1},
		}
		wantStatus := http.StatusOK
		requestBody, err := json.Marshal(reqBody)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, URL, bytes.NewBuffer(requestBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Basic "+basicAuth(conf.ConnectConfig.Login, conf.ConnectConfig.Password))
		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.NoError(t, err)
		assert.Equal(t, wantStatus, resp.StatusCode)

		ctx := context.Background()

		query := `SELECT status_id FROM orders WHERE order_id=$1;`
		var exists int64
		require.NotEqual(t, tdb, nil)
		err = tdb.QueryRow(ctx, query, 1).Scan(&exists)
		require.NoError(t, err)
		assert.Equal(t, exists, int64(2))
	})
}
