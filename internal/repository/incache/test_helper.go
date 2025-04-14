package incache

func setupCache() *Repo[string] {
	cache := NewRepository[string](3)
	cache.Storage[1] = &Node[string]{Key: 1, Value: "one"}
	cache.Storage[2] = &Node[string]{Key: 2, Value: "two"}
	cache.Storage[3] = &Node[string]{Key: 3, Value: "three"}

	cache.Head.Next = cache.Storage[1]
	cache.Storage[1].Prev = cache.Head
	cache.Storage[1].Next = cache.Storage[2]

	cache.Storage[2].Prev = cache.Storage[1]
	cache.Storage[2].Next = cache.Storage[3]

	cache.Storage[3].Prev = cache.Storage[2]
	cache.Storage[3].Next = cache.Tail

	cache.Tail.Prev = cache.Storage[3]
	return cache
}
