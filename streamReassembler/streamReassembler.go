package streamReassembler

import (
	"lab/stream"
	"strings"
)

type StreamReassemblerInterface interface {
	PushsubString(string, int, bool)
	StreamOut() stream.Stream
	UnassembledBytes() int
	Empty() bool
}

type StreamReassembler struct {
	outPut    *stream.Stream
	capacity  int
	eof       bool
	start     int
	startIdx  int
	buffer    []byte
	leftBytes int
	filled    []bool
}

var _ StreamReassemblerInterface = (*StreamReassembler)(nil)

func NewStreamReassembler(capacity int, output *stream.Stream) *StreamReassembler {
	output.Flush()
	return &StreamReassembler{
		outPut:    output,
		capacity:  capacity,
		eof:       false,
		start:     0,
		startIdx:  0,
		buffer:    make([]byte, capacity),
		leftBytes: 0,
		filled:    make([]bool, capacity),
	}
}

func (sr *StreamReassembler) StreamOut() stream.Stream {
	return *sr.outPut
}
func (sr *StreamReassembler) UnassembledBytes() int { return sr.leftBytes }
func (sr *StreamReassembler) Empty() bool           { return len(sr.buffer) == 0 }

func (sr *StreamReassembler) PushsubString(data string, idx int, eof bool) {
	//sr.outPut.Flush()
	from := 0

	if idx > sr.capacity || eof == true {
		sr.buffer = nil
		return
	}

	if idx < sr.startIdx {
		from = sr.startIdx - idx
	}
	size := min(len(data), sr.capacity-sr.outPut.BufferSize()-idx+sr.startIdx)
	for i := from; i < size; i++ {
		j := (i + sr.start + idx - sr.startIdx) % sr.capacity
		sr.buffer[j] = data[i]
		if !sr.filled[j] {
			sr.leftBytes++
		}
		sr.filled[j] = true
	}
	i := 0
	var segment strings.Builder
	for {
		j := (sr.start + i) % sr.capacity
		if !sr.filled[j] {
			break
		}
		sr.filled[j] = false
		sr.leftBytes--
		segment.WriteByte(sr.buffer[j])
		i++
	}
	sr.start = (sr.start + i) % sr.capacity
	sr.startIdx += i
	sr.outPut.Write(segment.String())
	sr.eof = sr.eof || (sr.eof && size == len(data))
	if sr.eof && sr.leftBytes == 0 {
		sr.outPut.EndInput()
	}
}
