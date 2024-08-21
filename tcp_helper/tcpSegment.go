package tcp_helper

import (
	"lab/stream"
	"lab/utils"
)

type TCPSegmentInterface interface {
	Parse(buffer utils.Buffer, dataGramLayerCheckSum uint32) utils.ParseResult
	Serialize(dataGramLayerCheckSum uint32) utils.BufferList
	GetHeader() TCPHeader[uint32]
	GetPayload() utils.Buffer
	LengthInSequenceSpace() int
}

type TCPSegment struct {
	header  TCPHeader[uint32]
	payload utils.Buffer
}

var _ TCPSegmentInterface = (*TCPSegment)(nil)

func NewTCPSegment(header TCPHeader[uint32]) *TCPSegment {
	return &TCPSegment{
		header:  header,
		payload: utils.Buffer{},
	}
}

func (s *TCPSegment) Parse(buffer utils.Buffer, dataGramLayerCheckSum uint32) utils.ParseResult {
	return 0
}

func (s *TCPSegment) Serialize(dataGramLayerCheckSum uint32) utils.BufferList {
	return utils.BufferList{}
}

func (s *TCPSegment) SetPaylaod(stream stream.Stream) {
	b := utils.ConvertStreamToBuffer(stream, stream.BufferSize())
	s.payload = *b
}

func (s *TCPSegment) GetHeader() TCPHeader[uint32] {
	return s.header
}

func (s *TCPSegment) GetPayload() utils.Buffer {
	return s.payload
}

func (s *TCPSegment) LengthInSequenceSpace() int {
	return 0
}
