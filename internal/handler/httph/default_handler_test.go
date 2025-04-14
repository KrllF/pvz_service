//go:build unit

package httph

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultHandler(t *testing.T) {
	t.Parallel()
	t.Run("all good", func(t *testing.T) {
		t.Parallel()
		serv := NewHandler(context.Background(),
			nil, nil, nil, nil)
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		serv.DefaultHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		assert.Equal(t, "Здравствуйте!", res.Body.String())
	})
}
