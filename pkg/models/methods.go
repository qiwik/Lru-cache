package models

type LRUCache interface {

}

func (c *Cache) Add(key string, value interface{}) bool {
	c.Lock()
	defer c.Unlock()
	if ok := c.validateItem(key); ok {
		return false
	}

	if c.Chain.Len() == c.Capacity {
		c.removeLast()
	}

	newItem := &Item{
		Key: key,
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

func (c *Cache) Remove(key string) bool {
	c.Lock()
	defer c.Unlock()
	if ok := c.validateItem(key); !ok {
		return false
	}

	item := c.Items[key]
	delete(c.Items, key)
	c.Chain.Remove(item)
	return true
}

func (c *Cache) validateItem(key string) bool {
	if _, ok := c.Items[key]; !ok {
		return false
	}
	return true
}

func (c *Cache) Get(key string) (value interface{}, ok bool) {
	c.Lock()
	defer c.Unlock()
	if ok := c.validateItem(key); !ok {
		return nil, false
	}

	value = c.Items[key].Value.(*Item).Value
	c.Chain.MoveToFront(c.Items[key])
	return value, true
}
