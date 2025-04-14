package incache

import "fmt"

func (c *Repo[T]) GetItem(key int64) (T, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.Storage[key]
	if !ok {
		var zero T
		return zero, fmt.Errorf("c.Storage[%d] отсутствует", key)
	}
	return item.Value, nil
}
