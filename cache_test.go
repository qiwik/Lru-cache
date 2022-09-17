package golru

import (
	lru "github.com/hashicorp/golang-lru"
	"github.com/stretchr/testify/require"
	"log"
	"strconv"
	"testing"
)

func TestInit(t *testing.T) {
	newCache, err := NewCache(2, 10)
	require.NoError(t, err)
	require.Equal(t, 2, newCache.capacity)
	require.Equal(t, 0, newCache.chain.Len())
	require.Equal(t, 0, len(newCache.items))
}

func TestInitSuccess(t *testing.T) {
	cache, err := NewCache(0, 10)
	require.Error(t, err)
	require.Nil(t, cache)
}

func TestAddPositive(t *testing.T) {
	newCache, err := NewCache(2, 10)
	require.NoError(t, err)

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
	newCache, err := NewCache(2, 10)
	require.NoError(t, err)

	answer := newCache.Add("test", 45)
	require.True(t, answer)
	fail := newCache.Add("test", "sos")
	require.False(t, fail)
	require.Equal(t, 1, newCache.chain.Len())
	require.Equal(t, 1, len(newCache.items))
}

func TestAddRemoveLast(t *testing.T) {
	newCache, err := NewCache(2, 10)
	require.NoError(t, err)

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
	newCache, err := NewCache(3, 10)
	require.NoError(t, err)

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
	newCache, err := NewCache(1, 10)
	require.NoError(t, err)

	answer := newCache.Add("test", 45)
	require.True(t, answer)

	value, ok := newCache.Get("testing")
	require.False(t, ok)
	require.Nil(t, value)
}

func TestRemovePositive(t *testing.T) {
	newCache, err := NewCache(1, 10)
	require.NoError(t, err)

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
	newCache, err := NewCache(1, 10)
	require.NoError(t, err)

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
	newCache, err := NewCache(3, 10)
	require.NoError(t, err)
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
	newCache, err := NewCache(1, 10)
	require.NoError(t, err)
	newCache.Add("test", 101)
	answer := newCache.ChangeValue("attempt", "pool")
	require.False(t, answer)
}

func TestRemoveAll(t *testing.T) {
	newCache, err := NewCache(3, 10)
	require.NoError(t, err)
	newCache.Add("test", 101)
	newCache.Add("tests", 102)
	newCache.Add("testing", 103)
	newCache.Clear()
	require.Equal(t, 0, newCache.Len())
}

func TestChangeCapacityToLarge(t *testing.T) {
	newCache, err := NewCache(1, 10)
	require.NoError(t, err)
	newCache.Add("test", 42)
	newCache.ChangeCapacity(2)
	newCache.Add("new", "testing")
	require.Equal(t, 2, newCache.Len())
}

func TestChangeCapacityToLess(t *testing.T) {
	newCache, err := NewCache(4, 10)
	require.NoError(t, err)
	newCache.Add("test1", 42)
	newCache.Add("test2", 43)
	newCache.Add("test3", 44)
	newCache.Add("test4", 45)
	newCache.ChangeCapacity(2)
	require.Equal(t, 2, newCache.Len())
}

func TestChangeCapacityToZero(t *testing.T) {
	newCache, err := NewCache(1, 10)
	require.NoError(t, err)
	newCache.Add("test", 42)
	newCache.ChangeCapacity(0)
	require.Equal(t, 1, newCache.Len())
}

func TestValues(t *testing.T) {
	newCache, err := NewCache(4, 10)
	require.NoError(t, err)
	newCache.Add("test1", 42)
	newCache.Add("test2", 43)
	newCache.Add("test3", 44)
	newCache.Add("test4", 45)
	values := newCache.Values()
	require.Len(t, values, 4)
}

func TestReflectKeys(t *testing.T) {
	newCache, err := NewCache(4, 10)
	require.NoError(t, err)
	newCache.Add("test1", 42)
	newCache.Add("test2", 43)
	newCache.Add("test3", 44)
	newCache.Add("test4", 45)
	keys := newCache.ReflectKeys()
	require.Len(t, keys, 4)
	require.IsType(t, "string", keys[0])
}

func TestKeys(t *testing.T) {
	newCache, err := NewCache(4, 10)
	require.NoError(t, err)
	newCache.Add("test1", 42)
	newCache.Add("test2", 43)
	newCache.Add("test3", 44)
	newCache.Add("test4", 45)
	keys := newCache.Keys()
	require.Len(t, keys, 4)
	require.IsType(t, "string", keys[0])
}

// Benchmarks

func BenchmarkReflectKeys(b *testing.B) {
	newCache, _ := NewCache(50, 10)
	for i := 0; i < 50; i++ {
		newCache.Add("test"+strconv.Itoa(i), 42)
	}
	for i := 0; i < b.N; i++ {
		newCache.ReflectKeys()
	}
}

func BenchmarkKeys(b *testing.B) {
	newCache, _ := NewCache(50, 10)
	for i := 0; i < 50; i++ {
		newCache.Add("test"+strconv.Itoa(i), 42)
	}
	for i := 0; i < b.N; i++ {
		newCache.Keys()
	}
}

func BenchmarkGolangLruKeys(b *testing.B) {
	cache, _ := lru.New(50)
	for i := 0; i < 50; i++ {
		cache.Add("test"+strconv.Itoa(i), 42)
	}
	for i := 0; i < b.N; i++ {
		cache.Keys()
	}
}

func BenchmarkGolangLruGet(b *testing.B) {
	cache, err := lru.New(100)
	if err != nil {
		log.Fatal("error")
	}
	cache.Add("key", struct {
		Slice     []string
		Integer   int
		NewStruct struct {
			Flo float64
		}
	}{
		Slice:     []string{"test", "testing", "golang"},
		Integer:   123456789987654321,
		NewStruct: struct{ Flo float64 }{Flo: 456215.12165468},
	})
	for i := 0; i < b.N; i++ {
		cache.Get("key")
	}
}

func BenchmarkGet(b *testing.B) {
	cache, err := NewCache(100, 10)
	if err != nil {
		log.Fatal("error")
	}
	cache.Add("key", struct {
		Slice     []string
		Integer   int
		NewStruct struct {
			Flo float64
		}
	}{
		Slice:     []string{"test", "testing", "golang"},
		Integer:   123456789987654321,
		NewStruct: struct{ Flo float64 }{Flo: 456215.12165468},
	})
	for i := 0; i < b.N; i++ {
		cache.Get("key")
	}
}

func BenchmarkGolangLruAdd(b *testing.B) {
	cache, err := lru.New(100)
	if err != nil {
		log.Fatal("error")
	}
	for i := 0; i < b.N; i++ {
		cache.Add("key", struct {
			Slice     []string
			Integer   int
			NewStruct struct {
				Flo float64
			}
		}{
			Slice:     []string{"test", "testing", "golang"},
			Integer:   123456789987654321,
			NewStruct: struct{ Flo float64 }{Flo: 456215.12165468},
		})
	}
}

func BenchmarkAdd(b *testing.B) {
	cache, err := NewCache(100, 10)
	if err != nil {
		log.Fatal("error")
	}
	for i := 0; i < b.N; i++ {
		cache.Add("key", struct {
			Slice     []string
			Integer   int
			NewStruct struct {
				Flo float64
			}
		}{
			Slice:     []string{"test", "testing", "golang"},
			Integer:   123456789987654321,
			NewStruct: struct{ Flo float64 }{Flo: 456215.12165468},
		})
	}
}

func BenchmarkGolangLruRemove(b *testing.B) {
	cache, err := lru.New(100)
	if err != nil {
		log.Fatal("error")
	}
	for i := 0; i < 100; i++ {
		cache.Add("key"+strconv.Itoa(i), struct {
			Slice     []string
			Integer   int
			NewStruct struct {
				Flo float64
			}
		}{
			Slice:     []string{"test", "testing", "golang"},
			Integer:   123456789987654321,
			NewStruct: struct{ Flo float64 }{Flo: 456215.12165468},
		})
	}
	for i := 0; i < b.N; i++ {
		cache.Remove("key")
	}
}

func BenchmarkRemove(b *testing.B) {
	cache, err := NewCache(100, 10)
	if err != nil {
		log.Fatal("error")
	}
	for i := 0; i < 100; i++ {
		cache.Add("key"+strconv.Itoa(i), struct {
			Slice     []string
			Integer   int
			NewStruct struct {
				Flo float64
			}
		}{
			Slice:     []string{"test", "testing", "golang"},
			Integer:   123456789987654321,
			NewStruct: struct{ Flo float64 }{Flo: 456215.12165468},
		})
	}
	for i := 0; i < b.N; i++ {
		cache.Remove("key")
	}
}
