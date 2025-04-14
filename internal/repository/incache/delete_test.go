package incache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteCache(t *testing.T) {
	t.Parallel()
	t.Run("all good", func(t *testing.T) {
		t.Parallel()
		cache := setupCache()

		err := cache.DeleteCache(2)
		require.NoError(t, err)

		_, ok := cache.Storage[2]
		assert.False(t, ok)

		assert.Equal(t, cache.Storage[3], cache.Storage[1].Next)
		assert.Equal(t, cache.Storage[1], cache.Storage[3].Prev)
	})

	t.Run("non-existent key", func(t *testing.T) {
		t.Parallel()
		cache := setupCache()

		err := cache.DeleteCache(999)
		assert.Error(t, err)
	})

	t.Run("deleting the first key", func(t *testing.T) {
		t.Parallel()
		cache := setupCache()

		err := cache.DeleteCache(1)
		require.NoError(t, err)

		_, ok := cache.Storage[1]
		assert.False(t, ok)

		assert.Equal(t, cache.Storage[2], cache.Head.Next)
	})

	t.Run("deleting the last key", func(t *testing.T) {
		t.Parallel()
		cache := setupCache()

		err := cache.DeleteCache(3)
		require.NoError(t, err)

		_, ok := cache.Storage[3]
		assert.False(t, ok)

		assert.Equal(t, cache.Tail.Prev, cache.Storage[2])
	})
}
