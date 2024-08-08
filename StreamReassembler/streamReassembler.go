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

func NewStreamReassembler(capacity int, output *stream.Stream) *StreamReassembler {
	return &StreamReassembler{
		unassembleStrs:     make(map[int]string),
		nextAssembledIdx:   0,
		unassebledBytesNum: 0,
		eofIdx:             -1,
		outPut:             output,
		capacity:           capacity,
	}
}

func (sr *StreamReassembler) PushSubString(data string, idx int, eof bool) {
	pos, f := MapFindUpperBoundIdx(sr.unassembleStrs, idx)
	log.Printf("pos: %d ,got = %v \n", pos, f)
	if !f {
		log.Printf("Error to find the target upperbound :%d\n", idx)
	}
	if pos != 0 {
		pos -= 1
	}
	newIdx := idx
	log.Println("new idx before change", newIdx)
	if f && pos <= idx {
		upIdx := pos
		log.Println("upIdx =", upIdx)
		if idx < upIdx+len(sr.unassembleStrs[upIdx]) {
			newIdx = upIdx + len(sr.unassembleStrs[upIdx])
			log.Println("new idx =", newIdx)
		} else {
			log.Println(upIdx + len(sr.unassembleStrs[upIdx]))
		}
	} else if idx < sr.nextAssembledIdx {
		newIdx = sr.nextAssembledIdx
		log.Println("new idx =", newIdx)
	}

	log.Println("new Idx after change", newIdx)

	dataStartPos := newIdx - idx
	dataSize := len(data) - dataStartPos
	log.Printf("dataStartPos = %d,dataSize = %d \n", dataStartPos, dataSize)
	for f && idx <= pos {
		pos += 1
		dataEndSize := newIdx + dataSize
		if pos < dataEndSize {
			if dataEndSize < pos+len(sr.unassembleStrs[pos]) {
				dataSize = pos - newIdx
				log.Println("new data size", dataSize)
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
	firstUnAcceptTableIdx := sr.nextAssembledIdx + sr.capacity - sr.outPut.BufferSize()

	if firstUnAcceptTableIdx <= newIdx {
		return
	}
	if dataSize > 0 {
		newData := data[:dataStartPos+dataSize]
		log.Printf("new data= %s,dataStartPos=%d,dataSize=%d\n", newData, dataStartPos, dataSize)
		if newIdx == sr.nextAssembledIdx {
			writeByte := sr.outPut.Write(newData)
			log.Println("write Byte", writeByte)
			sr.nextAssembledIdx += writeByte
			if writeByte < len(newData) {
				dataToStore := newData[writeByte : len(newData)-writeByte]
				sr.unassebledBytesNum += len(dataToStore)
				sr.unassembleStrs[sr.nextAssembledIdx] = dataToStore
				log.Println("unassemble string next idx", sr.unassembleStrs[sr.nextAssembledIdx])
			}
		} else {
			dataToStore := newData[0:len(newData)]
			sr.unassebledBytesNum += len(dataToStore)
			sr.unassembleStrs[newIdx] = dataToStore
			log.Println(sr.unassembleStrs[newIdx])
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
	var keys []int
	for k := range target {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	index := sort.Search(len(keys), func(i int) bool {
		return keys[i] > idx
	})
	if index < len(keys) {
		return keys[index], true
	}
	return -1, false
}
