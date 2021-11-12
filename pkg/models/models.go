package models

import (
	"container/list"
	"sync"
)

type Item struct {
	Key	string
	Value	interface{}
}

type Cache struct {
	sync.Mutex

	Capacity int
	Items	map[string]*list.Element
	Chain	*list.List
}

func NewLRUCache(n int) *Cache {
	return &Cache{
		Capacity: n,
		Items:    make(map[string]*list.Element),
		Chain:    list.New(),
	}
}