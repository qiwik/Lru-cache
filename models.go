package lru_cache

import (
	"container/list"
	"sync"
)

// Item is an element inside *list.Element of Cache
type Item struct {
	Key   string
	Value interface{}
}

type Cache struct {
	mu sync.Mutex

	Capacity int
	Items    map[string]*list.Element
	Chain    *list.List
}

// NewLRUCache create new implementation of Cache
func NewLRUCache(n int) *Cache {
	return &Cache{
		Capacity: n,
		Items:    make(map[string]*list.Element, n),
		Chain:    list.New(),
	}
}

// Add returns false if current key already exists, and true if key doesn't exist and new item was added to cache
func (c *Cache) Add(key string, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.validateItem(key)
	if ok {
		return false
	}

	if c.Chain.Len() == c.Capacity {
		c.removeLast()
	}

	newItem := &Item{
		Key:   key,
		Value: value,
	}
	newElement := c.Chain.PushFront(newItem)
	c.Items[newItem.Key] = newElement
	return true
}

func (c *Cache) removeLast() {
	currentElement := c.Chain.Back()
	item := c.Chain.Remove(currentElement).(*Item)
	delete(c.Items, item.Key)
}

// Remove returns false if current key doesn't exist, and true if removing was successful
func (c *Cache) Remove(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	element, ok := c.validateItem(key)
	if !ok {
		return false
	}

	delete(c.Items, key)
	c.Chain.Remove(element)
	return true
}

func (c *Cache) validateItem(key string) (*list.Element, bool) {
	if element, ok := c.Items[key]; !ok {
		return nil, false
	} else {
		return element, true
	}
}

// Get func returns a value with true if such element exist with current key, else returns nil and false
func (c *Cache) Get(key string) (value interface{}, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	element, ok := c.validateItem(key)
	if !ok {
		return nil, false
	}

	value = element.Value.(*Item).Value
	c.Chain.MoveToFront(element)
	return value, true
}
