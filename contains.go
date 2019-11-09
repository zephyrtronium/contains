// Package contains implements reusable, lightweight, efficient sets.
//
// The focus is on usefulness rather than set theory.
//
// Note that while this package provides two distinct set types, their APIs
// disagree in favor of efficiency.
//
package contains

import "math/bits"

const (
	// minDiff is the minimum number of bits a new key must add to be added
	// into an existing list.
	minDiff = 2
	// defCap is the capacity for new lists.
	defCap = 8
)

// A Set is a collection of sparse keys. The zero value is ready to use.
type Set struct {
	filters []uintptr
	keys    [][]uintptr
}

// Add adds the key to the set. Returns true if the key is new or false if the
// key was already present.
func (s *Set) Add(key uintptr) bool {
	r := filter(key)
	for k, f := range s.filters {
		if f&r == r {
			for _, v := range s.keys[k] {
				if v == key {
					return false
				}
			}
			// If the key is already present in a filter but not in the
			// associated list, we should add it to that list, so that further
			// checks will find it there.
			s.keys[k] = append(s.keys[k], key)
			return true
		}
	}
	k := len(s.filters) - 1
	// We want the new key to add at least minDiff bits to the filter. If it
	// won't, create a new list.
	if k >= 0 && bits.OnesCount64(uint64(r&^s.filters[k])) >= minDiff {
		s.filters[k] |= r
		s.keys[k] = append(s.keys[k], key)
	} else {
		s.filters = append(s.filters, r)
		// If we've previously reset, we might have extra lists available.
		if k+1 < cap(s.keys) {
			s.keys = s.keys[:k+2]
		} else {
			s.keys = append(s.keys, make([]uintptr, 0, defCap))
		}
		s.keys[k+1] = append(s.keys[k+1], key)
	}
	return true
}

// Contains returns true if key exists in the set.
func (s *Set) Contains(key uintptr) bool {
	r := filter(key)
	for k, f := range s.filters {
		if f&r == r {
			for _, v := range s.keys[k] {
				if v == key {
					return true
				}
			}
		}
	}
	return false
}

// Keys returns a slice containing all keys in the set. Returns nil if the set
// is empty.
func (s *Set) Keys() []uintptr {
	var r []uintptr
	for _, l := range s.keys {
		for _, v := range l {
			r = append(r, v)
		}
	}
	return r
}

// Reset removes all objects from the set. Reusing the set after calling Reset
// allows the previously allocated memory to be reused.
func (s *Set) Reset() {
	if s.filters != nil {
		s.filters = s.filters[:0]
		for k, v := range s.keys {
			s.keys[k] = v[:0]
		}
		// We don't have to resize s.keys itself because it's only ever read by
		// index and we check cap when adding to it. This lets us save a store
		// to memory, which has a significant impact in a tight loop.
	}
}

func filter(key uintptr) uintptr {
	if ^uintptr(0) != 0xffffffff {
		// 64-bit; use Knuth's MMIX LCG. We have to "convert" to uint64 because
		// these constants overflow uintptr on 32-bit, and the compiler doesn't
		// already know this branch is dead.
		return uintptr(6364136223846793005*uint64(key) + 1442695040888963407)
	}
	// 32-bit; use Numerical Recipes' LCG.
	return 1664525*key + 1013904223
}
