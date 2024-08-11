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

func TestBufferViewList(t *testing.T) {
	buf1 := NewBuffer("Hello")
	buf2 := NewBuffer("World")
	bl := NewBufferList(buf1)
	bl.Append(NewBufferList(buf2))
	bvl := NewBufferViewList(bl)
	if got := bvl.Size(); got != 10 {
		t.Errorf("BufferViewList.Size() = %d,want=%d", got, 10)
	}

	bvl.RemovePrefix(5)
	if got := bvl.Size(); got != 5 {
		t.Errorf("BufferViewList.Size() = %d,want=%d", got, 5)
	}

	iovecs := bvl.AsIOVecs()
	if len(iovecs) != 1 || string(iovecs[0]) != "World" {
		t.Errorf("BufferViewList.AsIOVecs() = %v,want =%v", iovecs, "World")
	}
}
