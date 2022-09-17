package golru

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	fromSecond      = 1000000000
	fromMillisecond = 1000000
)

var (
	ErrZeroTTL = errors.New("ttl should be greater than 0")
)

// todo: переделать на свою очередь

// Add returns false if current key already exists, and true if key doesn't exist and new item was added to cache.
// When a new element is added, it is placed at the top of the list, and if capacity is reached, the last element,
// which is also the most unpopular in the cache, is deleted
func (c *cache) Add(key string, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.validate(key); ok {
		return false
	}

	if c.chain.Len() == int(c.capacity) {
		c.removeLast()
	}

	newItem := &item{
		key:          key,
		value:        value,
		creationTime: time.Now(),
	}
	newElement := c.chain.PushFront(newItem)
	c.items[newItem.key] = newElement

	return true
}

// Get func returns a value with true if such element exist with current key, else returns nil and false. If an element
// exists, it is moved to the top of the list in the cache
func (c *cache) Get(key string) (interface{}, bool) {
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

// Remove returns false if current key doesn't exist, and true if removing from cache was successful
func (c *cache) Remove(key string) bool {
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

// ChangeValue allows you to change the value of a key that already exists in the cache. If there is no such key in
// the cache, the function returns false. If the value has changed, the element is sent to the top of the cache list
func (c *cache) ChangeValue(key string, newValue interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	element, ok := c.validate(key)
	if !ok {
		return false
	}

	element.Value.(*item).value = newValue
	element.Value.(*item).creationTime = time.Now()
	c.chain.MoveToFront(element)
	c.items[element.Value.(*item).key] = element

	return true
}

// Clear completely clears the cache
func (c *cache) Clear() {
	for c.chain.Len() > 0 {
		c.removeLast()
	}
}

// Len allows you to find out the fullness of the cache
func (c *cache) Len() int {
	return c.chain.Len()
}

// ChangeCapacity allows you to dynamically change the cache capacity. The new value must not be less than one. If
// the new capacity is less than the previous one, then the last elements in the list are deleted up to the desired
// parameter value
func (c *cache) ChangeCapacity(newCap uint32) {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch {
	case newCap <= 0:
		return
	case newCap >= c.capacity:
		c.capacity = newCap
		return
	default:
		c.capacity = newCap
		for c.Len() > int(newCap) {
			c.removeLast()
		}
	}
}

// Keys returns a slice of the keys that exist in the cache by simply traversing all the keys. Works faster than
// a function with reflection
func (c *cache) Keys() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	keys := make([]string, 0, len(c.items))
	for key := range c.items {
		keys = append(keys, key)
	}

	return keys
}

// ReflectKeys returns a slice of keys existing in the cache using reflection. It works 3-4 times slower than the Keys
// function, but is left for variability
func (c *cache) ReflectKeys() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	keysValues := reflect.ValueOf(c.items).MapKeys()
	keys := make([]string, 0, len(c.items))
	for i := range keysValues {
		keys = append(keys, keysValues[i].String())
	}

	return keys
}

// Values returns a slice of all existing element values in the cache
func (c *cache) Values() []interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	values := make([]interface{}, 0, len(c.items))
	for _, value := range c.items {
		values = append(values, value)
	}

	return values
}

// Expire starts checking the cache for the existence of expired data. Returns error if ttl is zero
func (c *cache) Expire(ctx context.Context) error {
	if c.ttl == 0 {
		return ErrZeroTTL
	}

	c.inspect()

	ticker := time.NewTicker(toNanosecond(float64(c.ttl)) * time.Nanosecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				c.inspect()
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()

	return nil
}

// inspect crawls the linked list and deletes data whose lifetime has come to an end
func (c *cache) inspect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	current := c.chain.Front()

	for current != nil {
		val := current.Value.(*item)
		if time.Now().Sub(val.creationTime).Seconds() > float64(c.ttl) {
			removed := current
			current = current.Next()

			delete(c.items, val.key)
			c.chain.Remove(removed)

			continue
		}

		current = current.Next()
	}
}

// validate checks the existence of an element by the key, and if it does not exist, returns false, instead of an element
func (c *cache) validate(key string) (element *list.Element, ok bool) {
	if element, ok = c.items[key]; !ok {
		return nil, false
	}
	return element, true
}

// removeLast deletes the last element in the list
func (c *cache) removeLast() {
	currentElement := c.chain.Back()
	last := c.chain.Remove(currentElement).(*item)
	delete(c.items, last.key)
}

// toNanosecond is a converter for ttl to time.Duration
func toNanosecond(ttl float64) time.Duration {
	ttlStr := fmt.Sprint(ttl)
	splitted := strings.Split(ttlStr, ".")

	f, _ := strconv.ParseInt(splitted[0], 10, 64)

	if len(splitted) != 1 {
		for len(splitted[1]) < 3 {
			splitted[1] += "0"
		}
		s, _ := strconv.ParseInt(splitted[1][:3], 10, 64)

		return time.Duration(f*fromSecond + s*fromMillisecond)
	}

	return time.Duration(f * fromSecond)
}
