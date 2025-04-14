package incache

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

func (c *Repo[T]) ListItems(ctx context.Context) ([]T, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ListItems from cache-service")
	defer span.Finish()
	var err error

	defer func() {
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("error", err.Error())
		}
	}()

	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]T, 0, len(c.Storage))

	current := c.Head.Next
	for current != c.Tail {
		result = append(result, current.Value)
		current = current.Next
	}

	return result, nil
}
