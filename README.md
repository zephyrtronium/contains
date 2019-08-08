# contains
Package contains implements a reusable set.

This is primarily intended to enable a fast cycle-avoiding graph traversal,
because `map[interface{}]struct{}` is slow. Typically the keys come from e.g.
`reflect.ValueOf(x).Pointer()`.
