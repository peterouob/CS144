package StreamReassembler

import "lab/stream"

type StreamReassemblerInterface interface {
	PushSubString(string, uint64, bool)
	StreamOut() stream.Stream
	UnassembledBytes() int
	Empty() bool
}

type StreamReassembler struct {
	unassembleStrs     map[int]string
	nextAssembledIdx   int
	unassebledBytesNum int
	eofIdx             int
	outPut             stream.Stream
	capacity           int
}

var _ StreamReassemblerInterface = (*StreamReassembler)(nil)

func New(next, num, eof int, output stream.Stream, capacity int) *StreamReassembler {
	return &StreamReassembler{
		unassembleStrs:     make(map[int]string),
		nextAssembledIdx:   next,
		unassebledBytesNum: num,
		eofIdx:             eof,
		outPut:             output,
		capacity:           capacity,
	}
}

func (sr *StreamReassembler) PushSubString(data string, idx uint64, eof bool) {}
func (sr *StreamReassembler) StreamOut() stream.Stream                        { return stream.Stream{} }
func (sr *StreamReassembler) UnassembledBytes() int                           { return sr.unassebledBytesNum }
func (sr *StreamReassembler) Empty() bool                                     { return sr.unassebledBytesNum == 0 }
