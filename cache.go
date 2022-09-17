package golru

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

type seconds float64
type CacheOption func(*cache)

var (
	ErrCacheCapacity = errors.New("capacity of the cache can not be less than 1")
)

// cache is a structure describing a temporary data store by key through a linked list. Mutex is present inside to
// implement the ability to work with the cache competitively.
// The number of entries is fixed. This implementation, rather than a memory limit, is chosen to work directly with
// storage objects. In addition to the list, there is a hash table with data inside the structure, which is used
// directly for accounting records.
// All fields are non-exportable, which allows you to work with the content through methods without having
// direct access to the cache
type cache struct {
	mu    sync.Mutex
	items map[string]*list.Element
	chain *list.List

	capacity uint32
	ttl      seconds
}

// NewCache create new implementation of lru cache. Capacity can't be less than one. If you set capacity to zero,
// for example, an assignment error will return
func NewCache(n uint32, opts ...CacheOption) (Cacher, error) {
	if n == 0 {
		return nil, ErrCacheCapacity
	}

	c := &cache{
		capacity: n,
		items:    make(map[string]*list.Element, n),
		chain:    list.New(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// item is an element inside *list.Element of cache with the key and value used by your program
type item struct {
	key   string
	value interface{}

	creationTime time.Time
}

func WithTTL(ttl seconds) CacheOption {
	return func(cache *cache) {
		cache.ttl = ttl
	}
}
