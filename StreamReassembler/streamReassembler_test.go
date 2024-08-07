package StreamReassembler

import "testing"

func TestMapFindUpperBoundIdx(t *testing.T) {
	tests := []struct {
		name   string
		target map[int]string
		idx    int
		want   int
		found  bool
	}{
		{"multiple elements", map[int]string{1: "one", 3: "three", 5: "five", 7: "seven"}, 4, 2, true},
		{"all elements less", map[int]string{1: "one", 3: "three", 5: "five"}, 7, -1, false},
		{"all elements more", map[int]string{2: "two", 4: "four", 6: "six"}, 1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := MapFindUpperBoundIdx(tt.target, tt.idx)
			if got != tt.want || found != tt.found {
				t.Errorf("MapFindUpperBoundIdx() = %d, found = %v, want %d, found %v", got, found, tt.want, tt.found)
			}
		})
	}
}
