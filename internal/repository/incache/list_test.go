package incache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListItems(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		cache := NewRepository[string](3)
		items, err := cache.ListItems(context.Background())
		assert.NoError(t, err)
		assert.Empty(t, items)
	})

	t.Run("items", func(t *testing.T) {
		cache := setupCache()

		items, err := cache.ListItems(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, []string{"one", "two", "three"}, items)
	})
}
