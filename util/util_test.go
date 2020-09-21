package util

import (
	"fyne.io/fyne"
	"testing"
)

func TestIntMax(t *testing.T) {
	tests := map[string]struct {
		numbers  []int
		expected int
	}{
		"normal numbers":   {numbers: []int{1, 2, 3, 4}, expected: 4},
		"negative numbers": {numbers: []int{-1, -2, -3, -4}, expected: -1},
		"mixed numbers":    {numbers: []int{1, 2, -3, -4}, expected: 2},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := IntMax(tc.numbers...)
			want := tc.expected
			if got != want {
				t.Fatalf("Expected %d, but got %d", want, got)
			}
		})
	}
}

func TestColumnMinSize(t *testing.T) {
	tests := map[string]struct {
		sizes    []fyne.Size
		expected fyne.Size
	}{
		"random sizes": {sizes: []fyne.Size{fyne.NewSize(5, 10), fyne.NewSize(20, 30)}, expected: fyne.NewSize(20, 40)},
		"zero sizes":   {sizes: []fyne.Size{fyne.NewSize(0, 0), fyne.NewSize(20, 30)}, expected: fyne.NewSize(20, 30)},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := ColumnMinSize(tc.sizes...)
			want := tc.expected
			if got.Width != want.Width || got.Height != want.Height {
				t.Fatalf("Expected height %d and width %d, got %#v", want.Height, want.Width, got)
			}
		})
	}
}

func TestInlineMinSize(t *testing.T) {
	tests := map[string]struct {
		sizes    []fyne.Size
		expected fyne.Size
	}{
		"random sizes": {sizes: []fyne.Size{fyne.NewSize(5, 10), fyne.NewSize(20, 30)}, expected: fyne.NewSize(25, 30)},
		"zero sizes":   {sizes: []fyne.Size{fyne.NewSize(0, 0), fyne.NewSize(20, 30)}, expected: fyne.NewSize(20, 30)},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := InlineMinSize(tc.sizes...)
			want := tc.expected
			if got.Width != want.Width || got.Height != want.Height {
				t.Fatalf("Expected height %d and width %d, got %#v", want.Height, want.Width, got)
			}
		})
	}
}
