package utils

import "testing"

func TestBufferList(t *testing.T) {

	buf1 := NewBuffer("Hello")
	buf2 := NewBuffer("World!")
	bl := NewBufferList(buf1)
	bl.Append(NewBufferList(buf2))
	if got := bl.Size(); got != 11 {
		t.Errorf("BufferList.Size() = %v, want %v", got, 11)
	}

	buffer, err := bl.ListToBuffer()
	if err != nil {
		t.Errorf("BufferList.ToBuffer() returned an error: %v", err)
	}
	if got := buffer.Str(); got != "HelloWorld!" {
		t.Errorf("BufferList.ToBuffer() = %v, want %v", got, "HelloWorld!")
	}

	bl.RemovePrefix(6)
	if got := bl.Size(); got != 5 {
		t.Errorf("BufferList.RemovePrefix() = %v, want %v", got, 5)
	}

	concatenated := bl.Concatenate()
	if concatenated != "orld!" {
		t.Errorf("BufferList.Concatenate() = %v, want %v", concatenated, "orld!")
	}
}
