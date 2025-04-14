package incache

import (
	"fmt"
)

func (c *Repo[T]) Put(key int64, value T) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Storage[key]; !ok {
		if err := c.addCache(key, value); err != nil {
			return fmt.Errorf("c.addCache: %w", err)
		}
		return nil
	}

	c.Storage[key].Value = value
	c.swapFirst(c.Storage[key])

	return nil
}

func (c *Repo[T]) addCache(key int64, value T) error {
	for int64(len(c.Storage)) >= c.Capacity {
		lastNode := c.Tail.Prev

		lastNode.Prev.Next = c.Tail
		c.Tail.Prev = lastNode.Prev

		delete(c.Storage, lastNode.Key)
	}

	newNode := NewNode(key, value)
	next := c.Head.Next
	c.Head.Next = newNode
	newNode.Prev = c.Head
	next.Prev = newNode
	newNode.Next = next

	c.Storage[key] = newNode

	return nil
}

func (c *Repo[T]) swapFirst(node *Node[T]) {
	if node == c.Head.Next {
		return
	}

	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev

	node.Next = c.Head.Next
	node.Prev = c.Head
	c.Head.Next.Prev = node
	c.Head.Next = node
}
