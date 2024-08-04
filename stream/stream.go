package stream

import "fmt"

type DequeFunc interface {
	Empty() bool
	PushFront(byte)
	PushBack(byte)
	PopFront() (byte, bool)
	PopBack() (byte, bool)
	StringItem() []string
}

type Deque struct {
	item []byte
}

var _ DequeFunc = (*Deque)(nil)

func (d *Deque) Empty() bool {
	return len(d.item) == 0
}

func (d *Deque) PushFront(item byte) {
	d.item = append([]byte{item}, d.item...)
}

func (d *Deque) PushBack(item byte) {
	d.item = append(d.item, item)
}

func (d *Deque) PopFront() (byte, bool) {
	if len(d.item) == 0 {
		return ' ', false
	}
	frontEle := d.item[0]
	d.item = d.item[1:]
	return frontEle, true
}

func (d *Deque) PopBack() (byte, bool) {
	if len(d.item) == 0 {
		return ' ', false
	}
	rearEle := d.item[len(d.item)-1]
	d.item = d.item[:len(d.item)-1]
	return rearEle, true
}

func (d *Deque) StringItem() []string {
	var tmp []string
	for _, v := range d.item {
		tmp = append(tmp, string(v))
	}
	return tmp
}

func newDeque(length int) *Deque {
	return &Deque{
		item: make([]byte, length),
	}
}

type StreamInterface interface {
	Write(string) int
	RemainingCapacity() int
	EndInput()
	SetError()
	PeekOutput(int) string
	PopOutPut(int)
	Read(int) string
	InputEnded() bool
	Errors() bool
	BufferSize() int
	BufferEmpty() bool
	EOF() bool
	BytesWritten() int
	BytesRead() int
}

type Stream struct {
	q            Deque
	capacitySize int
	writtenSize  int
	readSize     int
	endInput     bool
	error        bool
}

var _ StreamInterface = (*Stream)(nil)

func NewStream(q Deque, capacitySize, writtenSize, readSize int, endInput, error bool) *Stream {
	return &Stream{
		q:            q,
		capacitySize: capacitySize,
		writtenSize:  writtenSize,
		readSize:     readSize,
		endInput:     endInput,
		error:        error,
	}
}

func (stream *Stream) Write(data string) int {
	if stream.error || stream.endInput {
		return 0
	}
	writeSize := min(len(data), stream.RemainingCapacity())
	if writeSize == 0 {
		return 0
	}
	stream.writtenSize += writeSize
	for i := 0; i < writeSize; i++ {
		stream.q.PushBack(data[i])
	}
	return writeSize
}

func (stream *Stream) PeekOutput(length int) string {
	popSize := min(length, len(stream.q.item))
	return fmt.Sprintf("%s", stream.q.item[:popSize])
}

func (stream *Stream) PopOutPut(length int) {
	popSize := min(length, len(stream.q.item))
	stream.readSize += length
	for i := 0; i < popSize; i++ {
		stream.q.PopFront()
	}
}

func (stream *Stream) Read(length int) string {
	data := stream.PeekOutput(length)
	stream.PopOutPut(length)
	return data
}

func (stream *Stream) EndInput() {
	stream.endInput = true
}

func (stream *Stream) InputEnded() bool {
	return stream.endInput
}

func (stream *Stream) BufferSize() int {
	return len(stream.q.item)
}

func (stream *Stream) BufferEmpty() bool {
	return stream.q.Empty()
}

func (stream *Stream) EOF() bool {
	return stream.endInput && stream.q.Empty()
}

func (stream *Stream) BytesWritten() int {
	return stream.writtenSize
}

func (stream *Stream) BytesRead() int {
	return stream.readSize
}

func (stream *Stream) RemainingCapacity() int {
	return stream.capacitySize - len(stream.q.item)
}

func (stream *Stream) SetError() {
	stream.error = true
}

func (stream *Stream) Errors() bool {
	return stream.error
}
