package golru

import "context"

type Cacher interface {
	Expire(ctx context.Context) error

	Editor
	Informer
	Changer
}

type Editor interface {
	Add(key string, value interface{}) bool
	Get(key string) (interface{}, bool)
	Remove(key string) bool
	Clear()
}

type Changer interface {
	ChangeValue(key string, newValue interface{}) bool
	ChangeCapacity(newCap uint32)
}

type Informer interface {
	Len() int
	Keys() []string
	ReflectKeys() []string
	Values() []interface{}
}
