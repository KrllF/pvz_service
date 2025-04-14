package incache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPut(t *testing.T) {
	t.Parallel()

	t.Run("add new item", func(t *testing.T) {
		cache := setupCache()

		err := cache.Put(4, "four")
		assert.NoError(t, err)

		val, ok := cache.Storage[4]
		assert.True(t, ok)
		assert.Equal(t, "four", val.Value)

		assert.Equal(t, int64(4), cache.Head.Next.Key)
	})

	t.Run("update existing item", func(t *testing.T) {
		cache := setupCache()

		err := cache.Put(2, "new-two")
		assert.NoError(t, err)

		val, ok := cache.Storage[2]
		assert.True(t, ok)
		assert.Equal(t, "new-two", val.Value)

		assert.Equal(t, int64(2), cache.Head.Next.Key)
	})

	t.Run("exceed capacity", func(t *testing.T) {
		cache := setupCache()

		err := cache.Put(4, "four")
		assert.NoError(t, err)

		_, ok := cache.Storage[3]
		assert.False(t, ok)

		assert.Equal(t, int64(4), cache.Head.Next.Key)
	})
}
