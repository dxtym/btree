package btree_test

import (
	"testing"

	"github.com/dxtym/btree"
	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	items := []struct {
		key   int
		value string
	}{
		{1, "one"},
		{2, "two"},
		{3, "three"},
		{4, "four"},
		{5, "five"},
		{6, "six"},
		{7, "seven"},
		{8, "eight"},
		{9, "nine"},
		{10, "ten"},
	}

	t.Run("integer", func(t *testing.T) {
		b := btree.NewBtree[int, string](3) // TODO: make it variable

		for _, it := range items {
			b.Insert(it.key, it.value)
		}

		got := b.Traverse()
		assert.Len(t, got, 10)

		for i, it := range items {
			assert.Equal(t, it.key, got[i].Key)
			assert.Equal(t, it.value, got[i].Value)
		}
	})

}

func TestSearch(t *testing.T) {
	items := []struct {
		key   int
		value string
	}{
		{1, "one"},
		{2, "two"},
		{3, "three"},
		{4, "four"},
		{5, "five"},
		{6, "six"},
		{7, "seven"},
		{8, "eight"},
		{9, "nine"},
		{10, "ten"},
	}

	t.Run("", func(t *testing.T) {
		b := btree.NewBtree[int, string](3)

		for _, it := range items {
			b.Insert(it.key, it.value)
		}

		for _, it := range items {
			value, err := b.Search(it.key)
			assert.NoError(t, err)
			assert.Equal(t, it.value, value)
		}
	})
}
