package utils

type BufferInterface interface {
	Str() string
	At(int) (byte, bool)
	Size() int
	Copy() string
	RemovePrefix(int)
}

type Buffer struct {
	storage        []byte
	startingOffset int
}

var _ BufferInterface = (*Buffer)(nil)

func NewBuffer(str string) *Buffer {
	return &Buffer{
		storage: []byte(str),
	}
}

func (b *Buffer) Str() string {
	if b.storage == nil {
		return ""
	}
	return string(b.storage[b.startingOffset:])
}

func (b *Buffer) At(n int) (byte, bool) {
	if n >= b.Size() {
		return 0, false
	}
	return b.storage[b.startingOffset+n], true
}

func (b *Buffer) Size() int {
	return len(b.storage) - b.startingOffset
}

func (b *Buffer) Copy() string {
	return b.Str()
}

func (b *Buffer) RemovePrefix(n int) {
	if n > b.Size() {
		panic("prefix out of range")
	}
	b.startingOffset += n
}
