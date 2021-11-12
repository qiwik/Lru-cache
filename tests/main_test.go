package tests

import (
	"github.com/qiwik/lru-cache/pkg/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Init(t *testing.T) {
	newCache := models.NewLRUCache(2)
	require.Equal(t, 2, newCache.Capacity)
	require.Equal(t, 0, newCache.Chain.Len())
	require.Equal(t, 0, len(newCache.Items))
}

func Test_Add_Positive(t *testing.T) {
	newCache := models.NewLRUCache(2)

	answer := newCache.Add("test", 45)
	require.True(t, answer)
	require.Equal(t, 1, newCache.Chain.Len())
	require.Equal(t, 1, len(newCache.Items))
	item := newCache.Chain.Front()
	require.Equal(t, 45, item.Value.(*models.Item).Value)
	require.Equal(t, item, newCache.Items["test"])

	fail := newCache.Add("test", "sos")
	require.False(t, fail)
	require.Equal(t, 1, newCache.Chain.Len())
	require.Equal(t, 1, len(newCache.Items))
}

func Test_Add_Negative(t *testing.T) {
	newCache := models.NewLRUCache(2)

	answer := newCache.Add("test", 45)
	require.True(t, answer)
	fail := newCache.Add("test", "sos")
	require.False(t, fail)
	require.Equal(t, 1, newCache.Chain.Len())
	require.Equal(t, 1, len(newCache.Items))
}

func Test_Add_RemoveLast(t *testing.T) {
	newCache := models.NewLRUCache(2)

	answer1 := newCache.Add("test", 45)
	require.True(t, answer1)
	answer2 := newCache.Add("testing", "sos")
	require.True(t, answer2)
	answer3 := newCache.Add("new", 101)
	require.True(t, answer3)

	frontItem := newCache.Chain.Front()
	require.Equal(t, 101, frontItem.Value.(*models.Item).Value)
	require.Equal(t, frontItem, newCache.Items["new"])
	backItem := newCache.Chain.Back()
	require.Equal(t, "sos", backItem.Value.(*models.Item).Value)
	require.Equal(t, backItem, newCache.Items["testing"])
}

func Test_Get_Positive(t *testing.T) {
	newCache := models.NewLRUCache(3)

	answer1 := newCache.Add("test", 45)
	require.True(t, answer1)
	answer2 := newCache.Add("testing", "sos")
	require.True(t, answer2)
	answer3 := newCache.Add("new", 101)
	require.True(t, answer3)
	require.Equal(t, 3, newCache.Chain.Len())
	require.Equal(t, 3, len(newCache.Items))

	value, ok := newCache.Get("testing")
	require.True(t, ok)
	require.Equal(t, "sos", value)

	frontItem := newCache.Chain.Front()
	require.Equal(t, "sos", frontItem.Value.(*models.Item).Value)
}

func Test_Get_Negative(t *testing.T) {
	newCache := models.NewLRUCache(1)

	answer := newCache.Add("test", 45)
	require.True(t, answer)

	value, ok := newCache.Get("testing")
	require.False(t, ok)
	require.Nil(t, value)
}

func Test_Remove_Positive(t *testing.T) {
	newCache := models.NewLRUCache(1)

	answer := newCache.Add("test", 45)
	require.True(t, answer)
	require.Equal(t, 1, newCache.Chain.Len())
	require.Equal(t, 1, len(newCache.Items))

	ok := newCache.Remove("test")
	require.True(t, ok)
	require.Equal(t, 0, newCache.Chain.Len())
	require.Equal(t, 0, len(newCache.Items))
}

func Test_Remove_Negative(t *testing.T) {
	newCache := models.NewLRUCache(1)

	answer := newCache.Add("test", 45)
	require.True(t, answer)
	require.Equal(t, 1, newCache.Chain.Len())
	require.Equal(t, 1, len(newCache.Items))

	ok := newCache.Remove("tests")
	require.False(t, ok)
	require.Equal(t, 1, newCache.Chain.Len())
	require.Equal(t, 1, len(newCache.Items))
}