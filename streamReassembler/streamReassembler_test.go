package streamReassembler

import (
	"lab/stream"
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

func newStreamReassembler(next, num, eof int, data map[int]string, capacity int) *StreamReassembler {
	stream := stream.NewStream(stream.Deque{}, capacity, 0, 0, false, false)
	d := "hello peter"
	stream.Write(d)
	return NewStreamReassembler(capacity, stream)
}

func TestStreamReassembler_NormalInsertion(t *testing.T) {
	reassembler := newStreamReassembler(0, 0, -1, make(map[int]string), 20)
	t.Log(reassembler.outPut.Read(10))
	reassembler.unassembleStrs = map[int]string{
		1: "hello",
		2: "hey",
		3: "hi",
		6: "world",
	}
	reassembler.PushSubString("hello", 0, false)
	//reassembler.PushSubString("world", 5, true)

	out := reassembler.StreamOut()
	output := out.Read(5)
	expectedOutput := "hello"
	if output != expectedOutput {
		t.Errorf("Expected '%s' but got '%s'", expectedOutput, output)
	}

	if reassembler.UnassembledBytes() != 0 {
		t.Errorf("Expected 0 unassembled bytes but got %d", reassembler.UnassembledBytes())
	}
}

type MockStream struct {
	data []byte
}

func (s *MockStream) Write(data []byte) int {
	writeLen := len(data)
	if writeLen > len(s.data) {
		writeLen = len(s.data)
	}
	copy(s.data[:writeLen], data)
	return writeLen
}

func (s *MockStream) BufferSize() int {
	return len(s.data)
}

func (s *MockStream) EndInput() {
	// End of input logic here
}

func TestStreamReassembler_Normal(t *testing.T) {
	out := stream.NewStream(stream.Deque{}, 40, 0, 0, false, false)
	reassembler := NewStreamReassembler(20, out)
	reassembler.unassembleStrs = map[int]string{
		1: "hello world",
		2: "hey",
		3: "hi",
		6: "world",
	}
	reassembler.PushSubString("hello world", 0, false)
	//reassembler.PushSubString("peter", 3, false)
	reassembler.PushSubString("aaron", 1, false)

	expectedOutput := "d"
	gotOutput := out.Read(11)
	t.Log(gotOutput)
	gotOutput = out.Read(1)
	if gotOutput != expectedOutput {
		t.Errorf("Expected '%s' but got '%s'", expectedOutput, gotOutput)
	}
	if reassembler.unassebledBytesNum != 0 {
		t.Errorf("Expected unassembled bytes to be 0 but got %d", reassembler.unassebledBytesNum)
	}
}

func TestPushSubStringNormal(t *testing.T) {
	// Initialize the StreamReassembler
	out := stream.NewStream(stream.Deque{}, 40, 0, 0, false, false)
	sr := NewStreamReassembler(1024, out)
	// Mock data to push
	data := "hello world"
	sr.SetunassembleStrs(data, "hi")
	id := sr.GetTheIdxWithPayload(data)
	// Push substring
	sr.PushSubString(data, id, false)

	// Check the output buffer to see if data was written
	assembled := sr.outPut.Read(len(data)) // assuming outPut is accessible for testing
	if assembled != data {
		t.Errorf("need=%s,got=%s", data, assembled)
	}
}
