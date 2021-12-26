package lru_cache

import "container/list"

func (c *Cache) Len() int {
	return c.chain.Len()
}

func (c *Cache) Clear() {
	for c.chain.Len() > 0 {
		c.removeLast()
	}
}

func (c *Cache) ChangeCapacity(newCap int) {
	switch {
	case newCap <= 0:
		return
	case newCap >= c.capacity:
		c.capacity = newCap
		return
	default:
		c.capacity = newCap
		for c.Len() > newCap {
			c.removeLast()
		}
	}
}

func (c *Cache) ChangeValue(key string, newValue interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	element, ok := c.validate(key)
	if !ok {
		return false
	}

	element.Value.(*item).value = newValue
	c.chain.MoveToFront(element)
	return true
}

// Add returns false if current key already exists, and true if key doesn't exist and new item was added to cache
func (c *Cache) Add(key string, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.validate(key); ok {
		return false
	}

	if c.chain.Len() == c.capacity {
		c.removeLast()
	}

	newItem := &item{
		key:   key,
		value: value,
	}
	newElement := c.chain.PushFront(newItem)
	c.items[newItem.key] = newElement
	return true
}

func (c *Cache) removeLast() {
	currentElement := c.chain.Back()
	last := c.chain.Remove(currentElement).(*item)
	delete(c.items, last.key)
}

// Remove returns false if current key doesn't exist, and true if removing was successful
func (c *Cache) Remove(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	element, ok := c.validate(key)
	if !ok {
		return false
	}

	delete(c.items, key)
	c.chain.Remove(element)
	return true
}

func (c *Cache) validate(key string) (element *list.Element, ok bool) {
	if element, ok = c.items[key]; !ok {
		return nil, false
	}
	return element, true
}

// Get func returns a value with true if such element exist with current key, else returns nil and false
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	element, ok := c.validate(key)
	if !ok {
		return nil, false
	}

	value := element.Value.(*item).value
	c.chain.MoveToFront(element)
	return value, true
}
