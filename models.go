package lru_cache

import (
	"container/list"
	"sync"
)

type Cache struct {
	mu sync.Mutex

	capacity int
	items    map[string]*list.Element
	chain    *list.List
}

// item is an element inside *list.Element of Cache
type item struct {
	key   string
	value interface{}
}

// NewCache create new implementation of lru Cache
func NewCache(n int) *Cache {
	return &Cache{
		capacity: n,
		items:    make(map[string]*list.Element, n),
		chain:    list.New(),
	}
}
