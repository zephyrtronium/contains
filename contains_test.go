package contains

import (
	"fmt"
	"math/rand"
	"testing"
)

const testN = 1 << 16
const testLoops = 4

// TestContains tests that a Set does not contain values before adding and
// does contain them after adding.
func TestContains(t *testing.T) {
	v := make([]uintptr, testN)
	for i := range v {
		v[i] = uintptr(i)
	}
	s := Set{}
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

// TestAdd tests that a Set properly adds, rejects, and remembers values.
func TestAdd(t *testing.T) {
	v := make([]uintptr, testN)
	for i := range v {
		v[i] = uintptr(i)
	}
	s := Set{}
	for _, x := range v {
		if !s.Add(x) {
			t.Errorf("set failed to add key %d", x)
		}
		if !s.Contains(x) {
			t.Errorf("set lacks key %d", x)
		}
		if s.Add(x) {
			t.Errorf("set double-added key %d", x)
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
			if s.Add(x) {
				t.Errorf("set double-added key %d post-pop", x)
			}
		}
	}
}

// TestReset tests that a set contains no keys after resetting.
func TestReset(t *testing.T) {
	v := make([]uintptr, testN)
	for i := range v {
		v[i] = uintptr(i)
	}
	s := Set{}
	for _, x := range v {
		s.Add(x)
	}
	s.Reset()
	for _, x := range v {
		if !s.Add(x) {
			t.Errorf("couldn't readd key %d", x)
		}
	}
}

// BenchmarkContains benchmarks finding keys in a Set.
func BenchmarkContains(b *testing.B) {
	cases := []int{1 << 0, 1 << 2, 1 << 3, 1 << 6, 1 << 12, 1 << 16}
	for _, n := range cases {
		b.Run(fmt.Sprint(n), mkBContains(n))
	}
}

func mkBContains(n int) func(b *testing.B) {
	return func(b *testing.B) {
		s := Set{}
		for i := 0; i < n; i++ {
			s.Add(uintptr(i))
		}
		var v []int
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			k := i % n
			if k == 0 {
				v = rand.Perm(n)
			}
			if !s.Contains(uintptr(v[k])) {
				b.Errorf("set lost key %d", uintptr(v[k]))
			}
		}
	}
}

// BenchmarkAdd benchmarks adding keys to a Set.
func BenchmarkAdd(b *testing.B) {
	cases := []int{1 << 2, 1 << 3, 1 << 6, 1 << 12, 1 << 16}
	mod := []int{1 << 1, 1 << 3, 1 << 6, 1 << 12, 1 << 16}
	for _, n := range cases {
		for _, m := range mod {
			if m > n {
				break
			}
			b.Run(fmt.Sprintf("Realloc%d%%%d", n, m), mkBAddRealloc(n, m))
			b.Run(fmt.Sprintf("Reset%d%%%d", n, m), mkBAddReset(n, m))
		}
	}
}

func mkBAddRealloc(n, mod int) func(b *testing.B) {
	return func(b *testing.B) {
		v := make([]uintptr, n)
		for i := 0; i < n; i++ {
			v[i] = uintptr(i % mod)
		}
		s := Set{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h := i % n
			if h == 0 {
				// Test hangs if I pause the timer for this????
				for j := n - 1; j > 0; j-- {
					k := rand.Intn(j + 1)
					v[j], v[k] = v[k], v[j]
				}
				s = Set{}
			}
			s.Add(uintptr(h))
		}
	}
}

func mkBAddReset(n, mod int) func(b *testing.B) {
	return func(b *testing.B) {
		v := make([]uintptr, n)
		for i := 0; i < n; i++ {
			v[i] = uintptr(i % mod)
		}
		s := Set{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h := i % n
			if h == 0 {
				for j := n - 1; j > 0; j-- {
					k := rand.Intn(j + 1)
					v[j], v[k] = v[k], v[j]
				}
				s.Reset()
			}
			s.Add(uintptr(h))
		}
	}
}
