package StreamReassembler

import (
	"lab/stream"
	"log"
	"sort"
)

type StreamReassemblerInterface interface {
	PushSubString(string, int, bool)
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

func (sr *StreamReassembler) PushSubString(data string, idx int, eof bool) {
	pos, f := MapFindUpperBoundIdx(sr.unassembleStrs, idx)
	if !f {
		log.Panicf("Error to find the target upperbound :%d\n", idx)
	}
	if pos != 0 {
		pos -= 1
	}

}
func (sr *StreamReassembler) StreamOut() stream.Stream { return stream.Stream{} }
func (sr *StreamReassembler) UnassembledBytes() int    { return sr.unassebledBytesNum }
func (sr *StreamReassembler) Empty() bool              { return sr.unassebledBytesNum == 0 }
func MapFindUpperBoundIdx(target map[int]string, idx int) (int, bool) {
	if idx < 0 {
		return -1, false
	}
	var tmp = make([]int, 0, len(target))
	for i, _ := range target {
		tmp = append(tmp, i)
	}
	sort.Ints(tmp)
	index := sort.Search(len(tmp), func(i int) bool {
		return tmp[i] > idx
	})
	if index < len(tmp) {
		return index, true
	}
	return -1, false
}
