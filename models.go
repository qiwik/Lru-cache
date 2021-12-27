package golru

import (
	"container/list"
	"errors"
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
func NewCache(n int) (*Cache, error) {
	if n <= 0 {
		return nil, errors.New("capacity of the cache can not be less than 1")
	}
	return &Cache{
		capacity: n,
		items:    make(map[string]*list.Element, n),
		chain:    list.New(),
	}, nil
}
