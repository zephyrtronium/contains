# contains
[![GoDoc](https://godoc.org/github.com/zephyrtronium/contains?status.svg)](https://godoc.org/github.com/zephyrtronium/contains)

Package contains implements a reusable set.

This is primarily intended to enable a fast cycle-avoiding graph traversal,
because `map[interface{}]struct{}` is slow. Operations like union,
intersection, and symmetric difference are not provided, but can be
implemented. Typically the keys come from e.g. `reflect.ValueOf(x).Pointer()`.
