package incache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetItem(t *testing.T) {
	t.Parallel()
	t.Run("all good", func(t *testing.T) {
		t.Parallel()
		cache := setupCache()
		val, err := cache.GetItem(1)
		require.NoError(t, err)
		assert.Equal(t, "one", val)
	})

	t.Run("non-existent key", func(t *testing.T) {
		t.Parallel()
		cache := setupCache()
		_, err := cache.GetItem(999)
		assert.Error(t, err)
	})
}
