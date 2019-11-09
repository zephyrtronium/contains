package contains

import "math/bits"

// Dense is a set of integers backed by a single contiguous block of
// memory beginning at key 0. The zero value is ready to use.
type Dense struct {
	v []uint
}

const wordSize = bits.UintSize

// Add adds a key to the set. Unlike Set.Add, this returns nothing.
func (s *Dense) Add(key int) {
	w := key / wordSize
	if w >= len(s.v) {
		s.grow(w + 1)
	}
	s.v[w] |= 1 << uint(key%wordSize)
}

// Contains returns true if key exists in the set.
func (s *Dense) Contains(key int) bool {
	w := key / wordSize
	if w >= len(s.v) {
		return false
	}
	return s.v[w]&(1<<uint(key%wordSize)) != 0
}

// Reset removes all keys from the set. Reusing the set after calling Reset
// allows the previously allocated memory to be reused.
func (s *Dense) Reset() {
	if s.v != nil {
		s.v = s.v[:0]
	}
}

// Grow ensures that the backing array has sufficient space to hold the given
// key without needing to reallocate.
func (s *Dense) Grow(key int) {
	s.grow((key + wordSize - 1) / wordSize)
}

// grow is like Grow but for words instead of keys.
func (s *Dense) grow(w int) {
	switch {
	case w >= cap(s.v):
		// The current backing won't hold the new key. Reallocate.
		v := make([]uint, w)
		copy(v, s.v)
		s.v = v
	case w >= len(s.v):
		// The backing is large enough, but it's been reset. The previously
		// unused words might contain garbage from a previous use, so we need
		// to zero them.
		a := len(s.v)
		s.v = s.v[:w]
		for a < w {
			s.v[a] = 0
			a++
		}
	default:
		// Already large enough. Do nothing.
	}
}
