package contains

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
)

const testN = 1 << 14
const testLoops = 2

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

// TestKeys tests that a set returns the correct keys.
func TestKeys(t *testing.T) {
	v := make([]uintptr, testN)
	for i := range v {
		v[i] = uintptr(i)
	}
	// Test that Keys returns nil before adding any elements.
	s := Set{}
	t.Run("nothing_added", func(t *testing.T) {
		if u := s.Keys(); u != nil {
			t.Errorf("wrong keys: want nil, have %v", u)
		}
	})
	// Test that Keys does what we expect in the first place.
	for _, x := range v {
		s.Add(x)
	}
	u := s.Keys()
	// Sorting the keys places them in the same order as v.
	sort.Slice(u, func(i, j int) bool { return u[i] < u[j] })
	t.Run("long", func(t *testing.T) {
		if len(u) != len(v) {
			t.Fatalf("keys have incorrect length: want %d, have %d", len(v), len(u))
		}
		for i, x := range v {
			if x != u[i] {
				t.Errorf("incorrect key: want %d, have %d", x, u[i])
			}
		}
	})
	// Test that Keys works properly following Reset.
	s.Reset()
	s.Add(1)
	u = s.Keys()
	t.Run("reset", func(t *testing.T) {
		if len(u) != 1 {
			t.Fatalf("keys have incorrect length: want 1, have %d", len(u))
		}
		if u[0] != 1 {
			t.Errorf("wrong key: want 1, have %d", u[0])
		}
	})
	// Test taht Keys returns nil after resetting.
	s.Reset()
	t.Run("empty_reset", func(t *testing.T) {
		if u := s.Keys(); u != nil {
			t.Errorf("wrong keys: want nil, have %v", u)
		}
	})
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
		v := make([]uintptr, n)
		for i := 0; i < n; i++ {
			s.Add(uintptr(i))
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			k := i % n
			if !s.Contains(v[k]) {
				b.Errorf("set lost key %d", v[k])
			}
		}
	}
}

// BenchmarkAdd benchmarks adding keys to a Set.
func BenchmarkAdd(b *testing.B) {
	cases := []int{1 << 2, 1 << 3, 1 << 6, 1 << 12, 1 << 16}
	mod := []int{1 << 1, 1 << 3, 1 << 6, 1 << 12, 1 << 16}
	type benchcase struct {
		name string
		f    func(*testing.B)
	}
	var reallocs []benchcase
	var resets []benchcase
	for _, n := range cases {
		for _, m := range mod {
			if m > n {
				break
			}
			reallocs = append(reallocs, benchcase{fmt.Sprintf("Realloc_%dmod%d", n, m), mkBAddRealloc(n, m)})
			resets = append(resets, benchcase{fmt.Sprintf("Reset_%dmod%d", n, m), mkBAddReset(n, m)})
		}
	}
	for _, c := range reallocs {
		b.Run(c.name, c.f)
	}
	for _, c := range resets {
		b.Run(c.name, c.f)
	}
}

func mkBAddRealloc(n, mod int) func(*testing.B) {
	return func(b *testing.B) {
		v := make([]uintptr, n)
		for i := 0; i < n; i++ {
			v[i] = uintptr(i % mod)
		}
		s := Set{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			k := i % n
			if k == 0 {
				s = Set{}
			}
			s.Add(v[k])
		}
	}
}

func mkBAddReset(n, mod int) func(*testing.B) {
	return func(b *testing.B) {
		v := make([]uintptr, n)
		for i := 0; i < n; i++ {
			v[i] = uintptr(i % mod)
		}
		s := Set{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			k := i % n
			if k == 0 {
				s.Reset()
			}
			s.Add(v[k])
		}
	}
}

// ExampleSet shows an example of how to use a Set.
func ExampleSet() {
	s := Set{}
	fmt.Println("Contains 1:", s.Contains(1))
	fmt.Println("Added 1:", s.Add(1))
	fmt.Println("Contains 1:", s.Contains(1))
	fmt.Println("Added 1:", s.Add(1))
	s.Reset()
	fmt.Println("Added 1:", s.Add(1))
	// Output: Contains 1: false
	// Added 1: true
	// Contains 1: true
	// Added 1: false
	// Added 1: true
}
