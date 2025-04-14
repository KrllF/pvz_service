package incache

import (
	"fmt"
)

func (c *Repo[T]) DeleteCache(key int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, ok := c.Storage[key]
	if !ok {
		return fmt.Errorf("c.Storage[%d] не найден", key)
	}

	delete(c.Storage, key)
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev

	return nil
}
