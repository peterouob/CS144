package utils

import "testing"

func TestBuffer(t *testing.T) {
	buf := NewBuffer("Hello World")
	if got := buf.Str(); got != "Hello World" {
		t.Errorf("Buffer.Str() = %v, want %v", got, "Hello World")
	}

	if got, _ := buf.At(0); got != 'H' {
		t.Errorf("Buffer.At(0) = %v, want %v", got, 'H')
	}

	if got := buf.Size(); got != 11 {
		t.Errorf("Buffer.Size() = %v, want %v", got, 13)
	}

	buf.RemovePrefix(6)
	if got := buf.Str(); got != "World" {
		t.Errorf("Buffer.RemovePrefix() = %v, want %v", got, "World")
	}

	copied := buf.Copy()
	if copied != "World" {
		t.Errorf("Buffer.Copy() = %v, want %v", copied, "World")
	}
}
