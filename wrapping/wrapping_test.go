package wrapping

import (
	"testing"
)

func TestWrap(t *testing.T) {
	tests := []struct {
		name     string
		n        uint64
		isn      WrappingInt32
		expected uint32
	}{
		{"Basic wrap", 5, WrappingInt32{100}, 105},
		{"Wrap with overflow", 4294967295, WrappingInt32{1}, 0}, // 4294967295 is 2^32-1, which causes overflow
		{"Wrap with zero", 0, WrappingInt32{12345}, 12345},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WrappingInt32{}
			result := w.Wrap(tt.n, tt.isn)
			if result.rawValue != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result.rawValue)
			}
		})
	}
}

func TestUnWrap(t *testing.T) {
	tests := []struct {
		name        string
		n           WrappingInt32
		isn         WrappingInt32
		checkPoint  uint64
		expectedVal uint64
	}{
		{"Basic unwrap", WrappingInt32{105}, WrappingInt32{100}, 1000, 5},
		{"Unwrap with overflow", WrappingInt32{1}, WrappingInt32{4294967295}, 5000000000, 1164153216},
		{"Unwrap without overflow", WrappingInt32{0}, WrappingInt32{12345}, 12345, 4294954951},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WrappingInt32{}
			result := w.UnWrap(tt.n, tt.isn, tt.checkPoint)
			if result != tt.expectedVal {
				t.Errorf("expected %v, got %v", tt.expectedVal, result)
			}
		})
	}
}
