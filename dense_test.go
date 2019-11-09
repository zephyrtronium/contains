package contains

import (
	"fmt"
	"math/rand"
	"testing"
)

// TestDenseContains tests that a Dense set does not contain values before
// adding and does contain them after adding.
func TestDenseContains(t *testing.T) {
	v := make([]int, testN)
	for i := range v {
		v[i] = i
	}
	s := Dense{}
	for _, x := range v {
		if s.Contains(x) {
			t.Errorf("set has unexpected key %d", x)
		}
		s.Add(x)
		if !s.Contains(x) {
			t.Errorf("set lacks key %d", x)
		}
	}
	// Run extra loops to ensure behavior isn't different during and after
	// population and in random orders.
	for i := 0; i < testLoops; i++ {
		for j := testN - 1; j > 0; j-- {
			k := rand.Intn(i + 1)
			v[j], v[k] = v[k], v[j]
		}
		for _, x := range v {
			if !s.Contains(x) {
				t.Errorf("set lost key %d", x)
			}
		}
	}
}

// TestDenseAdd tests that a Dense set properly adds and remembers values.
func TestDenseAdd(t *testing.T) {
	s := Dense{}
	for x := 0; x < testN; x++ {
		s.Add(x)
		if !s.Contains(x) {
			t.Errorf("set lacks key %d", x)
		}
		s.Add(x)
		if !s.Contains(x) {
			t.Errorf("re-adding removed key %d", x)
		}
	}
}

// TestDenseReset tests that a Dense set contains no keys after resetting.
func TestDenseReset(t *testing.T) {
	s := Dense{}
	for x := 0; x < testN; x++ {
		s.Add(x)
	}
	s.Reset()
	for x := 0; x < testN; x++ {
		if s.Contains(x) {
			t.Errorf("set still contains key %d", x)
		}
	}
}

// TestDenseGrow tests that a Dense set contains the same keys after growing.
func TestDenseGrow(t *testing.T) {
	s := Dense{}
	for x := 0; x < testN/2; x++ {
		s.Add(x)
	}
	s.Grow(testN - 1)
	for x := 0; x < testN/2; x++ {
		if !s.Contains(x) {
			t.Errorf("set lost key %d", x)
		}
	}
	for x := testN / 2; x < testN; x++ {
		if s.Contains(x) {
			t.Errorf("set grew to have key %d", x)
		}
	}
	s.Grow(testN / 4)
	for x := 0; x < testN/2; x++ {
		if !s.Contains(x) {
			t.Errorf("set lost key %d after impotent grow", x)
		}
	}
	for x := testN / 2; x < testN; x++ {
		if s.Contains(x) {
			t.Errorf("set gained key %d after impotent grow", x)
		}
	}
	s.Reset()
	s.Grow(testN/2 - 1)
	for x := 0; x < testN; x++ {
		if s.Contains(x) {
			t.Errorf("set grew to have key %d after reset", x)
		}
	}
}

// ExampleDense shows an example of how to use a Dense set.
func ExampleDense() {
	s := Dense{}
	fmt.Println(s.Contains(1))
	s.Add(1)
	fmt.Println(s.Contains(1))
	s.Reset()
	fmt.Println(s.Contains(1))
	// Output: false
	// true
	// false
}
