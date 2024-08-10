package utils

import (
	"bytes"
)

type BufferListInterface interface {
	Buffers() []*Buffer
	Append(*BufferList)
	ListToBuffer() (*Buffer, error)
	RemovePrefix(int)
	Size() int
	Concatenate() string
}

type BufferList struct {
	buffers []*Buffer
}

var _ BufferListInterface = (*BufferList)(nil)

func NewBufferList(buffer *Buffer) *BufferList {
	return &BufferList{
		buffers: []*Buffer{buffer},
	}
}

func (b *BufferList) Buffers() []*Buffer {
	return b.buffers
}

func (b *BufferList) Append(other *BufferList) {
	b.buffers = append(b.Buffers(), other.buffers...)
}

func (b *BufferList) ListToBuffer() (*Buffer, error) {
	if len(b.buffers) == 0 {
		return NewBuffer(""), nil
	}
	if len(b.buffers) == 1 {
		return b.buffers[0], nil
	}

	var concatenated bytes.Buffer
	for _, buffer := range b.buffers {
		concatenated.Write([]byte(buffer.Str()))
	}
	return NewBuffer(concatenated.String()), nil
}

func (b *BufferList) RemovePrefix(n int) {
	for n > 0 && len(b.buffers) > 0 {
		buffer := b.buffers[0]
		bufferSize := buffer.Size()
		if n < bufferSize {
			buffer.RemovePrefix(n)
			n = 0
		} else {
			n -= bufferSize
			b.buffers = b.buffers[1:]
		}
	}
}

func (b *BufferList) Size() int {
	totalSize := 0
	for _, buffer := range b.buffers {
		totalSize += buffer.Size()
	}
	return totalSize
}

func (b *BufferList) Concatenate() string {
	var concatenated bytes.Buffer
	for _, buffer := range b.buffers {
		concatenated.Write([]byte(buffer.Str()))
	}
	return concatenated.String()
}
