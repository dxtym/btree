# Btree

An in-memory generic B-tree data structure in Go.

## Installation

```bash
go get github.com/dxtym/btree
```

## Example

```go
b, err := btree.New[int, string](3)
if err != nil {
    panic(err)
}

for range 100 {
    key := rand.Intn(100)
    b.Insert(key, strconv.Itoa(key))
}

value, err := b.Search(42)
if err != nil {
    panic(err)
}

if err := b.Remove(42); err != nil {
    panic(err)
}
```