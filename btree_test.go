package btree_test

import (
	"sort"
	"testing"

	"github.com/dxtym/btree"
	"github.com/stretchr/testify/assert"
)

type item struct {
	key   int
	value string
}

func TestBtree_Insert(t *testing.T) {
	tests := []struct {
		name   string
		order  int
		number int
		items  []item
	}{
		{
			name:   "simple",
			order:  3,
			number: 2,
			items: []item{
				{key: 1, value: "a"},
				{key: 2, value: "b"},
			},
		},
		{
			name:   "duplicate",
			order:  3,
			number: 2,
			items: []item{
				{key: 1, value: "a"},
				{key: 2, value: "b"},
				{key: 2, value: "c"},
			},
		},
		{
			name:   "one split",
			order:  3,
			number: 4,
			items: []item{
				{key: 1, value: "a"},
				{key: 2, value: "b"},
				{key: 3, value: "c"},
				{key: 4, value: "d"},
			},
		},
		{
			name:   "multiple split",
			order:  3,
			number: 7,
			items: []item{
				{key: 1, value: "a"},
				{key: 2, value: "b"},
				{key: 3, value: "c"},
				{key: 4, value: "d"},
				{key: 5, value: "e"},
				{key: 6, value: "f"},
				{key: 7, value: "g"},
			},
		},
		{
			name:   "random",
			order:  4,
			number: 9,
			items: []item{
				{key: 3, value: "c"},
				{key: 2, value: "b"},
				{key: 5, value: "e"},
				{key: 1, value: "a"},
				{key: 4, value: "d"},
				{key: 8, value: "h"},
				{key: 7, value: "g"},
				{key: 6, value: "f"},
				{key: 9, value: "i"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := btree.New[int, string](tt.order)

			for _, item := range tt.items {
				b.Insert(item.key, item.value)
			}

			items := make(map[int]string)
			for _, item := range tt.items {
				items[item.key] = item.value
			}

			result := b.Traverse()
			assert.Len(t, result, tt.number)
			assert.True(t, sort.SliceIsSorted(result, func(i, j int) bool {
				return result[i].Key < result[j].Key
			}))
			for _, res := range result {
				assert.Equal(t, items[res.Key], res.Value)
			}
		})
	}
}

func TestBtree_Search(t *testing.T) {
	tests := []struct {
		name    string
		order   int
		search  int
		items   []item
		wantErr error
	}{
		{
			name:   "exist",
			order:  3,
			search: 2,
			items: []item{
				{key: 1, value: "a"},
				{key: 2, value: "b"},
				{key: 3, value: "c"},
			},
			wantErr: nil,
		},
		{
			name:   "not exist",
			order:  3,
			search: 4,
			items: []item{
				{key: 1, value: "a"},
				{key: 2, value: "b"},
				{key: 3, value: "c"},
			},
			wantErr: btree.ErrKeyNotFound,
		},
		{
			name:   "deep",
			order:  3,
			search: 5,
			items: []item{
				{key: 1, value: "a"},
				{key: 2, value: "b"},
				{key: 3, value: "c"},
				{key: 4, value: "d"},
				{key: 5, value: "e"},
				{key: 6, value: "f"},
				{key: 7, value: "g"},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := btree.New[int, string](tt.order)

			for _, item := range tt.items {
				b.Insert(item.key, item.value)
			}

			value, err := b.Search(tt.search)
			assert.ErrorIs(t, err, tt.wantErr)

			for _, item := range tt.items {
				if item.key == tt.search {
					assert.Equal(t, item.value, value)
					break
				}
			}
		})
	}
}

func TestBtree_Remove(t *testing.T) {
	tests := []struct {
		name    string
		order   int
		key     int
		items   []item
		wantErr error
	}{
		{
			name:  "simple",
			order: 3,
			key:   2,
			items: []item{
				{key: 1, value: "a"},
				{key: 2, value: "b"},
				{key: 3, value: "c"},
			},
			wantErr: nil,
		},
		{
			name:  "not exist",
			order: 3,
			key:   4,
			items: []item{
				{key: 1, value: "a"},
				{key: 2, value: "b"},
				{key: 3, value: "c"},
			},
			wantErr: btree.ErrKeyNotFound,
		},
		{
			name:    "empty",
			order:   3,
			key:     1,
			items:   []item{},
			wantErr: btree.ErrTreeEmpty,
		},
		{
			name:  "leaf",
			order: 3,
			key:   3,
			items: []item{
				{key: 1, value: "a"},
				{key: 2, value: "b"},
				{key: 3, value: "c"},
				{key: 4, value: "d"},
				{key: 5, value: "e"},
				{key: 6, value: "f"},
				{key: 7, value: "g"},
			},
			wantErr: nil,
		},
		{
			name:  "internal",
			order: 3,
			key:   2,
			items: []item{
				{key: 1, value: "a"},
				{key: 2, value: "b"},
				{key: 3, value: "c"},
				{key: 4, value: "d"},
				{key: 5, value: "e"},
				{key: 6, value: "f"},
				{key: 7, value: "g"},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := btree.New[int, string](tt.order)

			for _, item := range tt.items {
				b.Insert(item.key, item.value)
			}

			err := b.Remove(tt.key)
			assert.ErrorIs(t, err, tt.wantErr) // TODO: fix dangling child

			result := b.Traverse()
			for _, res := range result {
				assert.NotEqual(t, tt.key, res.Key)
			}
		})
	}
}
