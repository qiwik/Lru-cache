package lru_cache

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInit(t *testing.T) {
	newCache := NewCache(2)
	require.Equal(t, 2, newCache.capacity)
	require.Equal(t, 0, newCache.chain.Len())
	require.Equal(t, 0, len(newCache.items))
}

func TestAddPositive(t *testing.T) {
	newCache := NewCache(2)

	answer := newCache.Add("test", 45)
	require.True(t, answer)
	require.Equal(t, 1, newCache.chain.Len())
	require.Equal(t, 1, len(newCache.items))
	elem := newCache.chain.Front()
	require.Equal(t, 45, elem.Value.(*item).value)
	require.Equal(t, elem, newCache.items["test"])

	fail := newCache.Add("test", "sos")
	require.False(t, fail)
	require.Equal(t, 1, newCache.chain.Len())
	require.Equal(t, 1, len(newCache.items))
}

func TestAddNegative(t *testing.T) {
	newCache := NewCache(2)

	answer := newCache.Add("test", 45)
	require.True(t, answer)
	fail := newCache.Add("test", "sos")
	require.False(t, fail)
	require.Equal(t, 1, newCache.chain.Len())
	require.Equal(t, 1, len(newCache.items))
}

func TestAddRemoveLast(t *testing.T) {
	newCache := NewCache(2)

	answer1 := newCache.Add("test", 45)
	require.True(t, answer1)
	answer2 := newCache.Add("testing", "sos")
	require.True(t, answer2)
	answer3 := newCache.Add("new", 101)
	require.True(t, answer3)

	frontItem := newCache.chain.Front()
	require.Equal(t, 101, frontItem.Value.(*item).value)
	require.Equal(t, frontItem, newCache.items["new"])
	backItem := newCache.chain.Back()
	require.Equal(t, "sos", backItem.Value.(*item).value)
	require.Equal(t, backItem, newCache.items["testing"])
}

func TestGetPositive(t *testing.T) {
	newCache := NewCache(3)

	answer1 := newCache.Add("test", 45)
	require.True(t, answer1)
	answer2 := newCache.Add("testing", "sos")
	require.True(t, answer2)
	answer3 := newCache.Add("new", 101)
	require.True(t, answer3)
	require.Equal(t, 3, newCache.chain.Len())
	require.Equal(t, 3, len(newCache.items))

	value, ok := newCache.Get("testing")
	require.True(t, ok)
	require.Equal(t, "sos", value)

	frontItem := newCache.chain.Front()
	require.Equal(t, "sos", frontItem.Value.(*item).value)
}

func TestGetNegative(t *testing.T) {
	newCache := NewCache(1)

	answer := newCache.Add("test", 45)
	require.True(t, answer)

	value, ok := newCache.Get("testing")
	require.False(t, ok)
	require.Nil(t, value)
}

func TestRemovePositive(t *testing.T) {
	newCache := NewCache(1)

	answer := newCache.Add("test", 45)
	require.True(t, answer)
	require.Equal(t, 1, newCache.chain.Len())
	require.Equal(t, 1, len(newCache.items))

	ok := newCache.Remove("test")
	require.True(t, ok)
	require.Equal(t, 0, newCache.chain.Len())
	require.Equal(t, 0, len(newCache.items))
}

func TestRemoveNegative(t *testing.T) {
	newCache := NewCache(1)

	answer := newCache.Add("test", 45)
	require.True(t, answer)
	require.Equal(t, 1, newCache.chain.Len())
	require.Equal(t, 1, len(newCache.items))

	ok := newCache.Remove("tests")
	require.False(t, ok)
	require.Equal(t, 1, newCache.chain.Len())
	require.Equal(t, 1, len(newCache.items))
}

func TestSetValueSuccess(t *testing.T) {
	newCache := NewCache(3)
	newCache.Add("test", 101)
	newCache.Add("tests", 102)
	newCache.Add("testing", 103)
	answer := newCache.ChangeValue("test", 0)
	require.True(t, answer)
	require.Equal(t, 0, newCache.items["test"].Value.(*item).value)
	require.Equal(t, 102, newCache.items["tests"].Value.(*item).value)
	require.Equal(t, 103, newCache.items["testing"].Value.(*item).value)

	frontItem := newCache.chain.Front()
	require.Equal(t, 0, frontItem.Value.(*item).value)
	backItem := newCache.chain.Back()
	require.Equal(t, 102, backItem.Value.(*item).value)
}

func TestSetValueInvalidKey(t *testing.T) {
	newCache := NewCache(1)
	newCache.Add("test", 101)
	answer := newCache.ChangeValue("attempt", "pool")
	require.False(t, answer)
}

func TestRemoveAll(t *testing.T) {
	newCache := NewCache(3)
	newCache.Add("test", 101)
	newCache.Add("tests", 102)
	newCache.Add("testing", 103)
	newCache.Clear()
	require.Equal(t, 0, newCache.Len())
}

func TestChangeCapacityToLarge(t *testing.T) {
	newCache := NewCache(1)
	newCache.Add("test", 42)
	newCache.ChangeCapacity(2)
	newCache.Add("new", "testing")
	require.Equal(t, 2, newCache.Len())
}

func TestChangeCapacityToLess(t *testing.T) {
	newCache := NewCache(4)
	newCache.Add("test1", 42)
	newCache.Add("test2", 43)
	newCache.Add("test3", 44)
	newCache.Add("test4", 45)
	newCache.ChangeCapacity(2)
	require.Equal(t, 2, newCache.Len())
}

func TestChangeCapacityToNull(t *testing.T) {
	newCache := NewCache(1)
	newCache.Add("test", 42)
	newCache.ChangeCapacity(0)
	require.Equal(t, 1, newCache.Len())
}
