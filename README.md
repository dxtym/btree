# Btree

[![license](https://img.shields.io/badge/License-MIT-blue)](#license)
[![codecov](https://codecov.io/gh/dxtym/btree/branch/main/graph/badge.svg)](https://codecov.io/gh/dxtym/btree)

An in-memory generic B-tree data structure in Go.

## Installation

```bash
go get github.com/dxtym/btree
```

## Features

- **Insert**: Add a key-value item
- **Search**: Find a value by key
- **Remove**: Delete an item by key
- **Traverse**: Get all items in order

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

if value, err := b.Search(42); err != nil {
    panic(err)
}

if err := b.Remove(42); err != nil {
    panic(err)
}
```

## Plans

- Make concurrent safe
- Add benchmarks tests
- Optimize memory usage

## License

MIT License. See [LICENSE](LICENSE) for details.