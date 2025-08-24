package btree

import "cmp"

type Btree[K cmp.Ordered, V any] struct {
	order int
	root  *node[K, V]
}

func New[K cmp.Ordered, V any](order int) (*Btree[K, V], error) {
	if order < 3 {
		return nil, ErrInvalidOrder
	}
	return &Btree[K, V]{
		order: order,
	}, nil
}

// Search finds the value of the given key.
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

// split performs split of the root node.
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

// Insert adds a new key-value pair to the tree.
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

// Traverse returns the in-order depth first search of the tree.
func (b *Btree[K, V]) Traverse() []*item[K, V] {
	if b.root != nil {
		return b.root.traverse()
	}

	return []*item[K, V]{}
}

// Remove deletes the key from the tree.
func (b *Btree[K, V]) Remove(key K) error {
	if b.root == nil {
		return ErrTreeEmpty
	}

	err, ok := b.root.remove(key)
	if ok {
		if !b.root.isEmpty() {
			return err
		}

		if b.root.isLeaf() {
			b.root = nil
		} else {
			b.root = b.root.children[0]
		}
	}

	return err
}
