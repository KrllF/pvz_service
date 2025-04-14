package incache

func (c *Repo[T]) UpdateAllCache(items []T, getKey func(T) int64) error {
	newStorage := make(map[int64]*Node[T])

	var prevNode *Node[T]
	for _, val := range items {
		key := getKey(val)
		node := NewNode(key, val)
		newStorage[key] = node

		if prevNode != nil {
			prevNode.Next = node
			node.Prev = prevNode
		}
		prevNode = node
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, oldNode := range c.Storage {
		oldNode.Prev = nil
		oldNode.Next = nil
	}

	if len(items) > 0 {
		firstKey := getKey(items[0])
		lastKey := getKey(items[len(items)-1])

		c.Head.Next = newStorage[firstKey]
		newStorage[firstKey].Prev = c.Head

		newStorage[lastKey].Next = c.Tail
		c.Tail.Prev = newStorage[lastKey]
	} else {
		c.Head.Next = c.Tail
		c.Tail.Prev = c.Head
	}

	c.Storage = newStorage

	return nil
}
