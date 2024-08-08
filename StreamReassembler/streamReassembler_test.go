package StreamReassembler

import (
	"lab/stream"
	"reflect"
	"testing"
)

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

func TestPushSubStringSubDataPosition(t *testing.T) {
	tests := []struct {
		name             string
		initialMap       map[int]string
		idx              int
		nextAssembledIdx int
		expectedNewIdx   int
		expectPanic      bool
	}{
		{"no substring", map[int]string{1: "one", 3: "three"}, 4, 2, 4, false},
		{"substring overlap", map[int]string{1: "one", 3: "three"}, 2, 0, 6, false},
		{"no substring and pos comparison", map[int]string{1: "one", 3: "three"}, 6, 2, 2, false},
		{"empty map", map[int]string{}, 4, 2, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectPanic {
						t.Errorf("unexpected panic: %v", r)
					}
				} else {
					if tt.expectPanic {
						t.Errorf("expected panic, but did not occur")
					}
				}
			}()
			sr := &StreamReassembler{
				unassembleStrs:   tt.initialMap,
				nextAssembledIdx: tt.nextAssembledIdx,
			}
			sr.PushSubString("data", tt.idx, false)
			// 在这里检查 newIdx 是否与预期匹配
		})
	}
}

func TestPushSubString(t *testing.T) {
	tests := []struct {
		name             string
		initialMap       map[int]string
		idx              int
		nextAssembledIdx int
		data             string
		eof              bool
		expectedNewIdx   int
		expectedMap      map[int]string
		expectedBytesNum int
	}{
		{
			name:             "No overlap",
			initialMap:       map[int]string{1: "one", 3: "three", 7: "four"},
			idx:              6,
			nextAssembledIdx: 4,
			data:             "",
			eof:              false,
			expectedNewIdx:   6,
			expectedMap:      map[int]string{1: "one", 3: "three", 7: "four"},
			expectedBytesNum: 8,
		},
		{
			name:             "Partial overlap",
			initialMap:       map[int]string{1: "one", 3: "three"},
			idx:              2,
			nextAssembledIdx: 0,
			data:             "xthree",
			eof:              false,
			expectedNewIdx:   6,
			expectedMap:      map[int]string{1: "one", 3: "three"},
			expectedBytesNum: 8,
		},
		{
			name:             "Complete overlap",
			initialMap:       map[int]string{1: "one", 3: "three"},
			idx:              3,
			nextAssembledIdx: 0,
			data:             "three",
			eof:              false,
			expectedNewIdx:   3,
			expectedMap:      map[int]string{1: "one", 3: "three"},
			expectedBytesNum: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := &StreamReassembler{
				unassembleStrs:     tt.initialMap,
				unassebledBytesNum: 8,
				nextAssembledIdx:   tt.nextAssembledIdx,
			}
			sr.PushSubString(tt.data, tt.idx, tt.eof)
			if sr.unassebledBytesNum != tt.expectedBytesNum {
				t.Errorf("expected bytes num: %d, got: %d", tt.expectedBytesNum, sr.unassebledBytesNum)
			}
			for k, v := range tt.expectedMap {
				if sr.unassembleStrs[k] != v {
					t.Errorf("expected map[%d]: %s, got: %s", k, v, sr.unassembleStrs[k])
				}
			}
		})
	}
}

func newStreamReassembler(m map[int]string, next, num, eof, capacity int) (*StreamReassembler, *stream.Stream) {
	strm := stream.NewStream(stream.Deque{}, capacity, 0, 0, false, false)
	reassembler := New(m, next, num, eof, strm, capacity)
	return reassembler, strm
}

func TestStreamReassembler(t *testing.T) {
	tests := []struct {
		initialMap map[int]string
		operations []struct {
			data  string
			index int
			eof   bool
		}
		expectedOutput string
		expectedMap    map[int]string // Expected internal map state
	}{
		{
			initialMap: map[int]string{},
			operations: []struct {
				data  string
				index int
				eof   bool
			}{
				{"hello", 0, false},
				{" world", 5, false},
			},
			expectedOutput: "hello world",
			expectedMap:    map[int]string{}, // After processing, the map should be empty
		},
		{
			initialMap: map[int]string{},
			operations: []struct {
				data  string
				index int
				eof   bool
			}{
				{"abc", 0, false},
				{"def", 3, false},
				{"ghi", 6, true},
			},
			expectedOutput: "abcdefghi",
			expectedMap:    map[int]string{}, // After processing, the map should be empty
		},
		{
			initialMap: map[int]string{},
			operations: []struct {
				data  string
				index int
				eof   bool
			}{
				{"part1", 0, false},
				{"part2", 5, false},
				{"part3", 10, true},
				{"extra", 15, false},
			},
			expectedOutput: "part1part2part3",
			expectedMap: map[int]string{
				15: "extra",
			}, // Extra data remains in the map
		},
		{
			initialMap: map[int]string{},
			operations: []struct {
				data  string
				index int
				eof   bool
			}{
				{"start", 0, false},
				{"end", 5, true},
				{"middle", 3, false},
			},
			expectedOutput: "startmiddleend",
			expectedMap: map[int]string{
				8: "end", // The 'end' part is still there due to overlapping
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.expectedOutput, func(t *testing.T) {
			reassembler, strm := newStreamReassembler(tt.initialMap, 0, 0, -1, 100) // Adjust capacity as needed

			for _, op := range tt.operations {
				reassembler.PushSubString(op.data, op.index, op.eof)
			}

			// Read the result from the stream
			output := strm.Read(len(tt.expectedOutput))

			// Check the result
			if output != tt.expectedOutput {
				t.Errorf("Expected %s but got %s", tt.expectedOutput, output)
			}

			// Check if the buffer size is as expected after operations
			if strm.BufferSize() != len(tt.expectedOutput) {
				t.Errorf("Buffer size is incorrect, got %d, expected %d", strm.BufferSize(), len(tt.expectedOutput))
			}

			// Check the internal map state of unassembleStrs
			if !reflect.DeepEqual(reassembler.unassembleStrs, tt.expectedMap) {
				t.Errorf("Expected unassembleStrs %v but got %v", tt.expectedMap, reassembler.unassembleStrs)
			}
		})
	}
}
