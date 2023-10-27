package main

import "testing"

func FuzzMe(f *testing.F) {
	// add some seed values that we want to ensure get tested
	f.Add(5, 10)
	f.Add(-5, -10)
	f.Add(0, 0)
	f.Add(1000000, -1000000)
	f.Fuzz(func(t *testing.T, x int, y int) {
		result := testAndFuzzMe(x, y)
		expected := x + y
		if result != expected {
			t.Errorf("Unexpected result for x = %d, y = %d: got %d, want %d", x, y, result, expected)
		}
	})
}

func TestMe(t *testing.T) {
	tests := []struct {
		x        int
		y        int
		expected int
	}{
		{5, 10, 15},
		{-5, -10, -15},
		{0, 0, 0},
		{1000000, -1000000, 0},
	}

	for _, test := range tests {
		result := testAndFuzzMe(test.x, test.y)
		if result != test.expected {
			t.Errorf("fuzzMe(%d, %d) = %d; want %d", test.x, test.y, result, test.expected)
		}
	}
}
