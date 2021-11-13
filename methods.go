package Lru_cache

type LRUCache interface {
	Add(key string, value interface{}) bool
	Remove(key string) bool
	Get(key string) (value interface{}, ok bool)
}
