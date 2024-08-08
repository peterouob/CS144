package StreamReassembler

import (
	"fmt"
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
	outPut             *stream.Stream
	capacity           int
}

var _ StreamReassemblerInterface = (*StreamReassembler)(nil)

func New(m map[int]string, next, num, eof int, output *stream.Stream, capacity int) *StreamReassembler {
	return &StreamReassembler{
		unassembleStrs:     m,
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
		log.Printf("Error to find the target upperbound :%d\n", idx)
	}
	if pos != 0 {
		pos -= 1
	}
	newIdx := idx
	if f && pos <= idx {
		upIdx := pos
		if idx < upIdx+len(sr.unassembleStrs[upIdx]) {
			newIdx = upIdx + len(sr.unassembleStrs[upIdx])
		}
	} else if idx < sr.nextAssembledIdx {
		newIdx = sr.nextAssembledIdx
	}

	dataStartPos := newIdx - idx
	dataSize := len(data) - dataStartPos

	for f && idx <= pos {
		pos += 1
		dataEndSize := newIdx + dataSize
		if pos < dataEndSize {
			if dataEndSize < pos+len(sr.unassembleStrs[pos]) {
				dataSize = pos - newIdx
				break
			} else {
				sr.unassebledBytesNum -= len(sr.unassembleStrs[pos])
				delete(sr.unassembleStrs, pos)
				pos += 1
				continue
			}
		} else {
			break
		}
	}
	//firstUnAcceptTableIdx := sr.nextAssembledIdx + sr.capacity - sr.outPut.BufferSize()
	firstUnAcceptTableIdx := sr.nextAssembledIdx + sr.capacity - 0

	if firstUnAcceptTableIdx <= newIdx {
		return
	}
	if dataSize > 0 {
		newData := data[dataStartPos : dataStartPos+dataSize]
		if newIdx == len(newData) {
			writeByte := sr.outPut.Write(newData)
			sr.nextAssembledIdx += writeByte
			if writeByte < len(newData) {
				dataToStore := newData[writeByte : len(newData)-writeByte]
				sr.unassebledBytesNum += len(dataToStore)
				sr.unassembleStrs[sr.nextAssembledIdx] = dataToStore
			}
		} else {
			dataToStore := newData[0:len(newData)]
			sr.unassebledBytesNum += len(dataToStore)
			sr.unassembleStrs[newIdx] = dataToStore
		}
	}
	for k, v := range sr.unassembleStrs {
		if sr.nextAssembledIdx <= k {
			log.Println(fmt.Sprintf("Assertion failed: nextAssembledIdx (%d) > iterFirst (%d)", sr.nextAssembledIdx, k))
		}
		if k == sr.nextAssembledIdx {
			writeNum := sr.outPut.Write(v)
			sr.nextAssembledIdx += writeNum
			if writeNum < len(v) {
				sr.unassebledBytesNum += len(v) - writeNum
				sr.unassembleStrs[sr.nextAssembledIdx] = v[:writeNum]
				sr.unassebledBytesNum -= len(v)
				delete(sr.unassembleStrs, k)
			}
			sr.unassebledBytesNum -= len(v)
			delete(sr.unassembleStrs, k)
		} else {
			break
		}
	}
	if eof {
		sr.eofIdx = idx + len(data)
	}
	if sr.eofIdx <= sr.nextAssembledIdx {
		sr.outPut.EndInput()
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
