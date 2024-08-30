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

func TestPushNotOverlap(t *testing.T) {
	q := stream.NewDeque(len("HelloWorld"))
	out := stream.NewStream(*q, 40, 10, 10)
	size := <-stream.Cap
	reassembler := NewStreamReassembler(size, out)

	reassembler.PushsubString("Hello", 0, false)
	expect := "Hello"
	if got := reassembler.outPut.ReadAll(); got != expect {
		t.Errorf("PushSubString() = %s, want = %s", got, expect)
	}

	reassembler.PushsubString("world", 5, false)
	expect = "Helloworld"
	if got := reassembler.outPut.ReadAll(); got != expect {
		t.Errorf("PushSubString() = %s, want = %s", got, expect)
	}
}

func TestPushOverlap(t *testing.T) {
	q := stream.NewDeque(len("HelloWorld"))
	out := stream.NewStream(*q, 40, 10, 10)
	size := <-stream.Cap
	reassembler := NewStreamReassembler(size, out)

	reassembler.PushsubString("Hello", 0, false)
	expect := "Hello"
	if got := reassembler.outPut.ReadAll(); got != expect {
		t.Errorf("PushSubString() = %s, want = %s", got, expect)
	}

	reassembler.PushsubString("world", 20, false)
	expect = "Hellold"
	if got := reassembler.outPut.ReadAll(); got != expect {
		t.Errorf("PushSubString() = %s, want = %s", got, expect)
	}
}

func TestOverCapacitySize(t *testing.T) {
	q := stream.NewDeque(len("HelloWorld"))
	out := stream.NewStream(*q, 40, 10, 10)
	size := <-stream.Cap
	reassembler := NewStreamReassembler(size, out)

	reassembler.PushsubString("Hello", 41, false)
	expect := ""
	if got := reassembler.outPut.ReadAll(); got != expect {
		t.Errorf("PushSubString() = %s, want = %s", got, expect)
	}
}

func TestEOF(t *testing.T) {
	q := stream.NewDeque(len("HelloWorld"))
	out := stream.NewStream(*q, 40, 10, 10)
	size := <-stream.Cap
	reassembler := NewStreamReassembler(size, out)
	reassembler.PushsubString("Hello", 0, true)
	expect := ""
	if got := reassembler.outPut.ReadAll(); got != expect && out.BufferEmpty() != true {
		t.Errorf("PushSubString() = %s, want = %s && Buffer want =%v,got = %v", got, expect, true, out.BufferEmpty())
	}
}
