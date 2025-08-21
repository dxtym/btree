package btree

import (
	"cmp"
	"errors"
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

// isEmpty checks whether the node is empty.
func (n *node[K, V]) isEmpty() bool {
	return n.numKeys == 0
}

// isLeaf checks whether the node is a leaf.
func (n *node[K, V]) isLeaf() bool {
	return n.numChildren == 0
}

// insertKey inserts a new key at the specified index.
func (n *node[K, V]) insertKey(index int, it *item[K, V]) {
	if index < n.numKeys {
		for i := n.numKeys; i > index; i-- {
			n.keys[i] = n.keys[i-1]
		}
	}
	n.keys[index] = it
	n.numKeys++
}

// insertChild inserts a new child at the specified index.
func (n *node[K, V]) insertChild(index int, child *node[K, V]) {
	if index < n.numChildren {
		for i := n.numChildren; i > index; i-- {
			n.children[i] = n.children[i-1]
		}
	}
	n.children[index] = child
	n.numChildren++
}

// removeKey removes the key at the specified index.
func (n *node[K, V]) removeKey(index int) {
	if index < n.numKeys-1 {
		for i := index; i < n.numKeys-1; i++ {
			n.keys[i] = n.keys[i+1]
		}
	}
	n.keys[n.numKeys-1] = nil
	n.numKeys--
}

// removeChild removes the child at the specified index.
func (n *node[K, V]) removeChild(index int) {
	if index < n.numChildren-1 {
		for i := index; i < n.numChildren-1; i++ {
			n.children[i] = n.children[i+1]
		}
	}
	n.children[n.numChildren-1] = nil
	n.numChildren--
}

// split separates the node to median and right.
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
	n.keys[mid] = nil

	diff := n.numKeys - mid - 1
	n.numKeys -= diff + 1 // NOTE: accounts for median key
	right.numKeys += diff

	if !n.isLeaf() {
		for i := mid + 1; i < n.numChildren; i++ {
			right.children[i-mid-1] = n.children[i]
			n.children[i] = nil
		}

		diff := n.numChildren - mid - 1
		n.numChildren -= diff
		right.numChildren += diff
	}

	return median, right
}

// search finds the position of the given key or its insertion index.
func (n *node[K, V]) search(key K) (int, bool) {
	low, high := 0, n.numKeys-1

	for low <= high {
		mid := low + (high-low)/2 // NOTE: avoids overflow
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

// insert adds a new item to the tree.
func (n *node[K, V]) insert(it *item[K, V]) bool {
	index, ok := n.search(it.Key)
	if ok {
		n.keys[index] = it
		return false
	}

	if n.isLeaf() {
		n.insertKey(index, it)
		return n.numKeys == n.order
	}

	if n.children[index].insert(it) {
		median, right := n.children[index].split()
		n.insertKey(index, median)
		n.insertChild(index+1, right)
	}

	return n.numKeys == n.order
}

// traverse returns in-order depth first search of the tree.
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

// borrowLeft borrows a key from the left sibling.
func (n *node[K, V]) borrowLeft(index int) {
	left, right := n.children[index-1], n.children[index]

	parent := n.keys[index-1]
	right.insertKey(0, parent)

	if !left.isLeaf() {
		right.insertChild(0, left.children[left.numChildren-1])
		left.removeChild(left.numChildren - 1)
	}

	n.keys[index-1] = left.keys[left.numKeys-1]
	left.removeKey(left.numKeys - 1)
}

// borrowRight borrows a key from the right sibling.
func (n *node[K, V]) borrowRight(index int) {
	left, right := n.children[index], n.children[index+1]

	parent := n.keys[index]
	left.insertKey(left.numKeys, parent)

	if !right.isLeaf() {
		left.insertChild(left.numChildren, right.children[0])
		right.removeChild(0)
	}

	n.keys[index] = right.keys[0]
	right.removeKey(0)
}

// merge merges the left node with its right sibling.
func (n *node[K, V]) merge(index int) {
	left, right := n.children[index], n.children[index+1]

	parent := n.keys[index]
	left.insertKey(left.numKeys, parent)

	for i := 0; i < right.numKeys; i++ {
		left.insertKey(left.numKeys, right.keys[i])
		right.keys[i] = nil
	}

	if !right.isLeaf() {
		for i := 0; i < right.numChildren; i++ {
			left.insertChild(left.numChildren, right.children[i])
			right.children[i] = nil
		}
	}

	n.removeKey(index)
	n.removeChild(index + 1)
}

// fill checks if the node has enough keys and fills it if necessary.
func (n *node[K, V]) fill(index int) {
	switch {
	case index > 0 && n.children[index-1].numKeys > n.order/2:
		n.borrowLeft(index)
	case index < n.numChildren && n.children[index+1].numKeys > n.order/2:
		n.borrowRight(index)
	default:
		n.merge(index)
	}
}

// getPredecessor retrieves the predecessor of the key at the specified index.
func (n *node[K, V]) getPredecessor(index int) *item[K, V] {
	curr := n.children[index]
	for !curr.isLeaf() {
		curr = curr.children[curr.numChildren-1]
	}
	pred := curr.keys[curr.numKeys-1]
	curr.removeKey(curr.numKeys - 1)
	return pred
}

// getSuccessor retrieves the successor of the key at the specified index.
func (n *node[K, V]) getSuccessor(index int) *item[K, V] {
	curr := n.children[index+1]
	for !curr.isLeaf() {
		curr = curr.children[0]
	}
	succ := curr.keys[0]
	curr.removeKey(0)
	return succ
}

// remove deletes the key from the tree.
func (n *node[K, V]) remove(key K) (error, bool) {
	index, ok := n.search(key)
	if ok {
		if n.isLeaf() {
			n.removeKey(index)
		} else {
			switch {
			case n.children[index].numKeys > n.order/2:
				pred := n.getPredecessor(index)
				n.keys[index] = pred
			case n.children[index+1].numKeys > n.order/2:
				succ := n.getSuccessor(index)
				n.keys[index] = succ
			default:
				n.merge(index)
			}
		}

		return nil, n.numKeys < n.order/2
	}

	if n.isLeaf() {
		return errors.New("key not found"), false
	}

	err, ok := n.children[index].remove(key)
	if ok {
		n.fill(index)
	}

	return err, n.numKeys < n.order/2
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
	return value, errors.New("key not found")
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
		return errors.New("tree is empty")
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
