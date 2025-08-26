package btree_test

import (
	"sync"
	"testing"

	"github.com/dxtym/btree"
	"github.com/stretchr/testify/assert"
)

func TestBtree_Insert(t *testing.T) {
	tests := []struct {
		name   string
		order  int
		number int
		items  map[int]string
	}{
		{
			name:   "Simple",
			order:  3,
			number: 2,
			items: map[int]string{
				1: "a",
				2: "b",
			},
		},
		{
			name:   "One Split",
			order:  3,
			number: 4,
			items: map[int]string{
				1: "a",
				2: "b",
				3: "c",
				4: "d",
			},
		},
		{
			name:   "Multiple Split",
			order:  3,
			number: 7,
			items: map[int]string{
				1: "a",
				2: "b",
				3: "c",
				4: "d",
				5: "e",
				6: "f",
				7: "g",
			},
		},
		{
			name:   "Random",
			order:  4,
			number: 9,
			items: map[int]string{
				3: "c",
				2: "b",
				5: "e",
				1: "a",
				4: "d",
				8: "h",
				7: "g",
				6: "f",
				9: "i",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := btree.New[int, string](tt.order)

			wg := sync.WaitGroup{}
			for key, value := range tt.items {
				wg.Add(1)
				go func() {
					b.Insert(key, value)
					wg.Done()
				}()
			}
			wg.Wait()

			items := b.Iter()
			for it := range items {
				assert.Equal(t, tt.items[it.Key], it.Value)
			}
		})
	}
}

func TestBtree_Search(t *testing.T) {
	tests := []struct {
		name    string
		order   int
		search  int
		items   map[int]string
		wantErr error
	}{
		{
			name:   "Exist",
			order:  3,
			search: 2,
			items: map[int]string{
				1: "a",
				2: "b",
				3: "c",
			},
			wantErr: nil,
		},
		{
			name:   "Not Exist",
			order:  3,
			search: 4,
			items: map[int]string{
				1: "a",
				2: "b",
				3: "c",
			},
			wantErr: btree.ErrKeyNotFound,
		},
		{
			name:   "Deep",
			order:  3,
			search: 5,
			items: map[int]string{
				1: "a",
				2: "b",
				3: "c",
				4: "d",
				5: "e",
				6: "f",
				7: "g",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := btree.New[int, string](tt.order)

			wg := sync.WaitGroup{}
			for key, value := range tt.items {
				wg.Add(1)
				go func() {
					b.Insert(key, value)
					wg.Done()
				}()
			}
			wg.Wait()

			found, err := b.Search(tt.search)
			assert.ErrorIs(t, err, tt.wantErr)

			for key, value := range tt.items {
				if key == tt.search {
					assert.Equal(t, value, found)
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
		items   map[int]string
		wantErr error
	}{
		{
			name:  "Simple",
			order: 3,
			key:   2,
			items: map[int]string{
				1: "a",
				2: "b",
				3: "c",
			},
			wantErr: nil,
		},
		{
			name:  "Not Exist",
			order: 3,
			key:   4,
			items: map[int]string{
				1: "a",
				2: "b",
				3: "c",
			},
			wantErr: btree.ErrKeyNotFound,
		},
		{
			name:  "Leaf",
			order: 3,
			key:   3,
			items: map[int]string{
				1: "a",
				2: "b",
				3: "c",
				4: "d",
				5: "e",
				6: "f",
				7: "g",
			},
			wantErr: nil,
		},
		{
			name:  "Internal",
			order: 3,
			key:   2,
			items: map[int]string{
				1: "a",
				2: "b",
				3: "c",
				4: "d",
				5: "e",
				6: "f",
				7: "g",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := btree.New[int, string](tt.order)

			wg := sync.WaitGroup{}
			for key, value := range tt.items {
				wg.Add(1)
				go func() {
					b.Insert(key, value)
					wg.Done()
				}()
			}
			wg.Wait()

			wg.Add(1)
			go func() {
				err := b.Remove(tt.key)
				assert.ErrorIs(t, err, tt.wantErr)
				wg.Done()
			}()
			wg.Wait()

			items := b.Iter()
			for it := range items {
				assert.NotEqual(t, tt.key, it.Key)
			}
		})
	}
}

func TestBtree_Iter(t *testing.T) {
	tests := []struct {
		name  string
		order int
		items map[int]string
	}{
		{
			name:  "Simple",
			order: 3,
			items: map[int]string{
				1: "a",
				2: "b",
				3: "c",
				4: "d",
				5: "e",
				6: "f",
				7: "g",
			},
		},
		{
			name:  "Reverse",
			order: 3,
			items: map[int]string{
				7: "g",
				6: "f",
				5: "e",
				4: "d",
				3: "c",
				2: "b",
				1: "a",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := btree.New[int, string](tt.order)

			wg := sync.WaitGroup{}
			for key, value := range tt.items {
				wg.Add(1)
				go func() {
					b.Insert(key, value)
					wg.Done()
				}()
			}
			wg.Wait()

			items := b.Iter()
			for it := range items {
				assert.Equal(t, tt.items[it.Key], it.Value)
			}
		})
	}
}
