package btree

import "errors"

var (
	ErrTreeEmpty    = errors.New("tree is empty")
	ErrKeyNotFound  = errors.New("key not found")
	ErrInvalidOrder = errors.New("order must be at least 3")
)
