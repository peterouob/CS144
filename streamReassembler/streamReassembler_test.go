package streamReassembler

import (
	"lab/stream"
	"testing"
)

func TestPush(t *testing.T) {
	q := stream.NewDeque(len("HelloWorld"))
	out := stream.NewStream(*q, 40, 10, 10)
	reassembler := NewStreamReassembler(20, out)
	reassembler.PushsubString("Hello", 0, false)
	expect := "Hello"
	if got := reassembler.outPut.ReadAll(); got != expect {
		t.Errorf("PushsubString() = %s,want=%s", got, expect)
	}
	reassembler.PushsubString("World", 5, false)
	expect = "HelloWorld"
	if got := reassembler.outPut.ReadAll(); got != expect {
		t.Errorf("PushsubString() = %s,want=%s", got, expect)
	}
}

func TestPushOverlap(t *testing.T) {
	q := stream.NewDeque(len("HelloWorld"))
	out := stream.NewStream(*q, 40, 10, 10)
	reassembler := NewStreamReassembler(20, out)

	// Push the first part
	reassembler.PushsubString("Hello", 0, false)
	expect := "Hello"
	if got := reassembler.outPut.ReadAll(); got != expect {
		t.Errorf("PushSubString() = %s, want = %s", got, expect)
	}

	// Push an overlapping part
	reassembler.PushsubString("World", 4, false)
	expect = "HelWorld"
	if got := reassembler.outPut.ReadAll(); got != expect {
		t.Errorf("PushSubString() = %s, want = %s", got, expect)
	}
}
