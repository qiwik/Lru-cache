package golru

import (
	"context"
	"log"
	"strconv"
	"testing"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	c, err := NewCache(2)
	require.NoError(t, err)

	tc := c.(*cache)
	require.Equal(t, uint32(2), tc.capacity)
	require.Equal(t, 0, tc.chain.Len())
	require.Equal(t, 0, len(tc.items))
}

func TestInitSuccess(t *testing.T) {
	c, err := NewCache(0)
	require.Error(t, err)
	require.Nil(t, c)
}

func TestAddPositive(t *testing.T) {
	c, err := NewCache(2)
	require.NoError(t, err)

	tc := c.(*cache)

	answer := c.Add("test", 45)
	require.True(t, answer)
	require.Equal(t, 1, tc.chain.Len())
	require.Equal(t, 1, len(tc.items))

	elem := tc.chain.Front()
	require.Equal(t, 45, elem.Value.(*item).value)
	require.Equal(t, elem, tc.items["test"])

	fail := c.Add("test", "sos")
	require.False(t, fail)
	require.Equal(t, 1, tc.chain.Len())
	require.Equal(t, 1, len(tc.items))
}

func TestAddNegative(t *testing.T) {
	c, err := NewCache(2)
	require.NoError(t, err)

	tc := c.(*cache)

	answer := c.Add("test", 45)
	require.True(t, answer)

	fail := c.Add("test", "sos")
	require.False(t, fail)
	require.Equal(t, 1, tc.chain.Len())
	require.Equal(t, 1, len(tc.items))
}

func TestAddRemoveLast(t *testing.T) {
	c, err := NewCache(2)
	require.NoError(t, err)

	tc := c.(*cache)

	answer1 := c.Add("test", 45)
	require.True(t, answer1)
	answer2 := c.Add("testing", "sos")
	require.True(t, answer2)
	answer3 := c.Add("new", 101)
	require.True(t, answer3)

	frontItem := tc.chain.Front()
	require.Equal(t, 101, frontItem.Value.(*item).value)
	require.Equal(t, frontItem, tc.items["new"])

	backItem := tc.chain.Back()
	require.Equal(t, "sos", backItem.Value.(*item).value)
	require.Equal(t, backItem, tc.items["testing"])
}

func TestGetPositive(t *testing.T) {
	c, err := NewCache(3)
	require.NoError(t, err)

	tc := c.(*cache)

	answer1 := c.Add("test", 45)
	require.True(t, answer1)
	answer2 := c.Add("testing", "sos")
	require.True(t, answer2)
	answer3 := c.Add("new", 101)
	require.True(t, answer3)
	require.Equal(t, 3, tc.chain.Len())
	require.Equal(t, 3, len(tc.items))

	value, ok := c.Get("testing")
	require.True(t, ok)
	require.Equal(t, "sos", value)

	frontItem := tc.chain.Front()
	require.Equal(t, "sos", frontItem.Value.(*item).value)
}

func TestGetNegative(t *testing.T) {
	c, err := NewCache(1)
	require.NoError(t, err)

	answer := c.Add("test", 45)
	require.True(t, answer)

	value, ok := c.Get("testing")
	require.False(t, ok)
	require.Nil(t, value)
}

func TestRemovePositive(t *testing.T) {
	c, err := NewCache(1)
	require.NoError(t, err)

	tc := c.(*cache)

	answer := c.Add("test", 45)
	require.True(t, answer)
	require.Equal(t, 1, tc.chain.Len())
	require.Equal(t, 1, len(tc.items))

	ok := c.Remove("test")
	require.True(t, ok)
	require.Equal(t, 0, tc.chain.Len())
	require.Equal(t, 0, len(tc.items))
}

func TestRemoveNegative(t *testing.T) {
	c, err := NewCache(1)
	require.NoError(t, err)

	tc := c.(*cache)

	answer := c.Add("test", 45)
	require.True(t, answer)
	require.Equal(t, 1, tc.chain.Len())
	require.Equal(t, 1, len(tc.items))

	ok := c.Remove("tests")
	require.False(t, ok)
	require.Equal(t, 1, tc.chain.Len())
	require.Equal(t, 1, len(tc.items))
}

func TestSetValueSuccess(t *testing.T) {
	c, err := NewCache(3)
	require.NoError(t, err)

	tc := c.(*cache)

	c.Add("test", 101)
	c.Add("tests", 102)
	c.Add("testing", 103)

	answer := c.ChangeValue("test", 0)
	require.True(t, answer)
	require.Equal(t, 0, tc.items["test"].Value.(*item).value)
	require.Equal(t, 102, tc.items["tests"].Value.(*item).value)
	require.Equal(t, 103, tc.items["testing"].Value.(*item).value)

	frontItem := tc.chain.Front()
	require.Equal(t, 0, frontItem.Value.(*item).value)
	backItem := tc.chain.Back()
	require.Equal(t, 102, backItem.Value.(*item).value)
}

func TestSetValueInvalidKey(t *testing.T) {
	c, err := NewCache(1)
	require.NoError(t, err)

	c.Add("test", 101)
	answer := c.ChangeValue("attempt", "pool")
	require.False(t, answer)
}

func TestRemoveAll(t *testing.T) {
	c, err := NewCache(3)
	require.NoError(t, err)

	c.Add("test", 101)
	c.Add("tests", 102)
	c.Add("testing", 103)
	c.Clear()
	require.Equal(t, 0, c.Len())
}

func TestChangeCapacityToLarge(t *testing.T) {
	c, err := NewCache(1)
	require.NoError(t, err)

	c.Add("test", 42)
	c.ChangeCapacity(2)
	c.Add("new", "testing")
	require.Equal(t, 2, c.Len())
}

func TestChangeCapacityToLess(t *testing.T) {
	c, err := NewCache(4)
	require.NoError(t, err)

	c.Add("test1", 42)
	c.Add("test2", 43)
	c.Add("test3", 44)
	c.Add("test4", 45)
	c.ChangeCapacity(2)
	require.Equal(t, 2, c.Len())
}

func TestChangeCapacityToZero(t *testing.T) {
	c, err := NewCache(1)
	require.NoError(t, err)

	c.Add("test", 42)
	c.ChangeCapacity(0)
	require.Equal(t, 1, c.Len())
}

func TestValues(t *testing.T) {
	c, err := NewCache(4)
	require.NoError(t, err)

	c.Add("test1", 42)
	c.Add("test2", 43)
	c.Add("test3", 44)
	c.Add("test4", 45)

	values := c.Values()
	require.Len(t, values, 4)
}

func TestReflectKeys(t *testing.T) {
	c, err := NewCache(4)
	require.NoError(t, err)

	c.Add("test1", 42)
	c.Add("test2", 43)
	c.Add("test3", 44)
	c.Add("test4", 45)

	keys := c.ReflectKeys()
	require.Len(t, keys, 4)
	require.IsType(t, "string", keys[0])
}

func TestKeys(t *testing.T) {
	c, err := NewCache(4)
	require.NoError(t, err)

	c.Add("test1", 42)
	c.Add("test2", 43)
	c.Add("test3", 44)
	c.Add("test4", 45)

	keys := c.Keys()
	require.Len(t, keys, 4)
	require.IsType(t, "string", keys[0])
}

func TestExpireFractional(t *testing.T) {
	c, err := NewCache(2, WithTTL(0.5))
	require.NoError(t, err)

	tc := c.(*cache)

	c.Add("test ttl", "ttl")
	require.Equal(t, tc.chain.Len(), 1)
	require.Len(t, tc.items, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = c.Expire(ctx)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)
	require.Equal(t, tc.chain.Len(), 0)
	require.Len(t, tc.items, 0)
}

func TestExpireInteger(t *testing.T) {
	c, err := NewCache(2, WithTTL(1))
	require.NoError(t, err)

	tc := c.(*cache)

	c.Add("test ttl", "ttl")
	c.Add("test ttl 2", "ttl")
	require.Equal(t, tc.chain.Len(), 2)
	require.Len(t, tc.items, 2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = c.Expire(ctx)
	require.NoError(t, err)

	time.Sleep(3 * time.Second)
	require.Equal(t, tc.chain.Len(), 0)
	require.Len(t, tc.items, 0)
}

func TestExpireZeroTTL(t *testing.T) {
	c, err := NewCache(2)
	require.NoError(t, err)

	c.Add("test ttl", "ttl")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = c.Expire(ctx)
	require.ErrorIs(t, err, ErrZeroTTL)
}

// Benchmarks

func BenchmarkReflectKeys(b *testing.B) {
	newCache, _ := NewCache(50)
	for i := 0; i < 50; i++ {
		newCache.Add("test"+strconv.Itoa(i), 42)
	}
	for i := 0; i < b.N; i++ {
		newCache.ReflectKeys()
	}
}

func BenchmarkKeys(b *testing.B) {
	newCache, _ := NewCache(50)
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
	cache, err := NewCache(100)
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
	cache, err := NewCache(100)
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
	cache, err := NewCache(100)
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
