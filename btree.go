package btree

import (
	"cmp"
)

type item[K cmp.Ordered, V any] struct {
	Key   K
	Value V
}

type node[K cmp.Ordered, V any] struct {
	order       int
	numKeys     int
	numChildren int
	keys        []*item[K, V]
	children    []*node[K, V]
}

func (n *node[K, V]) isLeaf() bool {
	return n.numChildren == 0
}

func (n *node[K, V]) split() (*item[K, V], *node[K, V]) {
	mid := n.numKeys / 2
	median := n.keys[mid]

	right := &node[K, V]{
		order:    n.order,
		keys:     make([]*item[K, V], n.order),
		children: make([]*node[K, V], n.order+1),
	}

	for i := mid + 1; i < n.numKeys; i++ {
		right.keys[i-mid-1] = n.keys[i]
		n.keys[i] = nil
	}
	n.numKeys -= mid + 1
	right.numKeys = n.numKeys
	n.keys[mid] = nil

	if !n.isLeaf() {
		for i := mid + 1; i < n.numChildren; i++ {
			right.children[i-mid-1] = n.children[i]
			n.children[i] = nil
		}
		n.numChildren -= mid + 1
		right.numChildren = n.numChildren
	}

	return median, right
}

func (n *node[K, V]) search(key K) (int, bool) {
	low, high := 0, n.numKeys-1

	for low <= high {
		mid := low + (high-low)/2
		current := n.keys[mid].Key

		if current > key {
			high = mid - 1
		} else if current < key {
			low = mid + 1
		} else {
			return mid, true
		}
	}

	return low, false
}

func (n *node[K, V]) insert(it *item[K, V]) bool {
	index, ok := n.search(it.Key)
	if ok {
		n.keys[index] = it
		return false
	}

	if n.isLeaf() {
		n.keys[index] = it
		n.numKeys++
		return n.numKeys == n.order
	}

	if n.children[index].insert(it) {
		median, right := n.children[index].split()
		n.keys[n.numKeys] = median
		n.children[n.numChildren] = right
		n.numKeys++
		n.numChildren++
	}

	return n.numKeys == n.order
}

func (n *node[K, V]) traverse() []*item[K, V] {
	if n.isLeaf() {
		items := make([]*item[K, V], n.numKeys)
		for i := 0; i < n.numKeys; i++ {
			items[i] = n.keys[i]
		}
		return items
	}

	items := make([]*item[K, V], 0, n.numKeys) // TODO: optimize capacity
	for i := 0; i < n.numChildren; i++ {
		items = append(items, n.children[i].traverse()...)
		if i < n.numKeys {
			items = append(items, n.keys[i])
		}
	}

	return items
}

type Btree[K cmp.Ordered, V any] struct {
	order int
	root  *node[K, V]
}

func NewBtree[K cmp.Ordered, V any](order int) *Btree[K, V] {
	return &Btree[K, V]{
		order: order,
	}
}

func (b *Btree[K, V]) Search(key K) (V, error) {
	for node := b.root; node != nil; {
		index, ok := node.search(key)
		if ok {
			return node.keys[index].Value, nil
		}

		node = node.children[index]
	}

	var value V
	return value, ErrKeyNotFound
}

func (b *Btree[K, V]) split() {
	root := &node[K, V]{
		order:    b.order,
		keys:     make([]*item[K, V], b.order),
		children: make([]*node[K, V], b.order+1),
	}

	median, right := b.root.split()
	root.keys[root.numKeys] = median
	root.children[root.numChildren] = b.root
	root.children[root.numChildren+1] = right
	root.numKeys++
	root.numChildren += 2

	b.root = root
}

func (b *Btree[K, V]) Insert(key K, value V) {
	it := &item[K, V]{
		Key:   key,
		Value: value,
	}

	if b.root == nil {
		b.root = &node[K, V]{
			order:    b.order,
			keys:     make([]*item[K, V], b.order),
			children: make([]*node[K, V], b.order+1),
		}
	}

	if b.root.insert(it) {
		b.split()
	}
}

func (b *Btree[K, V]) Traverse() []*item[K, V] {
	if b.root != nil {
		return b.root.traverse()
	}

	return []*item[K, V]{}
}
