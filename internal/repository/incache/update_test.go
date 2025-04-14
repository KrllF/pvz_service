package incache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getKey(value string) int64 {
	return int64(len(value))
}

func TestUpdateAllCache(t *testing.T) {
	t.Parallel()

	t.Run("update cache", func(t *testing.T) {
		cache := setupCache()

		newItems := []string{"a", "ab", "abc"}
		err := cache.UpdateAllCache(newItems, getKey)
		assert.NoError(t, err)

		_, ok := cache.Storage[0]
		assert.False(t, ok)

		assert.Equal(t, "a", cache.Storage[1].Value)
		assert.Equal(t, "ab", cache.Storage[2].Value)
		assert.Equal(t, "abc", cache.Storage[3].Value)

		assert.Equal(t, int64(1), cache.Head.Next.Key)
		assert.Equal(t, int64(2), cache.Head.Next.Next.Key)
	})

	t.Run("update cache with empty slice", func(t *testing.T) {
		cache := setupCache()

		err := cache.UpdateAllCache([]string{}, getKey)
		assert.NoError(t, err)

		assert.Empty(t, cache.Storage)
		assert.Equal(t, cache.Tail, cache.Head.Next)
		assert.Equal(t, cache.Head, cache.Tail.Prev)
	})
}
