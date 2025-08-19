package btree_test

import (
	"math/rand/v2"
	"sort"
	"strings"
	"testing"

	"github.com/dxtym/btree"
	"github.com/stretchr/testify/assert"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateRandomNumber(t testing.TB, n int) int {
	t.Helper()

	for {
		num := rand.IntN(n)
		if num >= 2 {
			return num
		}
	}
}

func generateRandomNumbers(t testing.TB, n int) []int {
	t.Helper()
	values := make([]int, n)
	for i := range values {
		values[i] = i + 1
	}
	rand.Shuffle(n, func(i, j int) {
		values[i], values[j] = values[j], values[i]
	})
	return values
}

func generateRandomStrings(t testing.TB, n, size int) []string {
	t.Helper()
	values := make([]string, n)
	for i := range values {
		var sb strings.Builder
		for range size {
			sb.WriteByte(letters[rand.IntN(len(letters))])
		}
		values[i] = sb.String()
	}
	return values
}

func TestBtree_InsertNumber(t *testing.T) {
	type testItem struct {
		key   int
		value string
	}

	order := generateRandomNumber(t, 10)
	keys := generateRandomNumbers(t, 100)
	values := generateRandomStrings(t, 100, 10)

	items := make([]testItem, 100)
	for i := range items {
		items[i] = testItem{
			key:   keys[i],
			value: values[i],
		}
	}

	b := btree.NewBtree[int, string](order)

	for _, it := range items {
		b.Insert(it.key, it.value)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].key < items[j].key
	})

	traversal := b.Traverse()
	assert.Len(t, traversal, 100)

	for i, tr := range traversal {
		assert.Equal(t, items[i].key, tr.Key)
		assert.Equal(t, items[i].value, tr.Value)
	}
}


func TestBtree_InsertString(t *testing.T) {
	type testItem struct {
		key   string
		value string
	}

	order := generateRandomNumber(t, 10)
	keys := generateRandomStrings(t, 100, 10)
	values := generateRandomStrings(t, 100, 10)

	items := make([]testItem, 100)
	for i := range items {
		items[i] = testItem{
			key:   keys[i],
			value: values[i],
		}
	}

	b := btree.NewBtree[string, string](order)

	for _, it := range items {
		b.Insert(it.key, it.value)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].key < items[j].key
	})

	traversal := b.Traverse()
	assert.Len(t, traversal, 100)

	for i, tr := range traversal {
		assert.Equal(t, items[i].key, tr.Key)
		assert.Equal(t, items[i].value, tr.Value)
	}
}
