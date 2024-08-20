package tcp_helper

import "lab/utils"

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

func NewTCPSegment() *TCPSegment {
	return &TCPSegment{
		header:  TCPHeader[uint32]{},
		payload: utils.Buffer{},
	}
}

func (s *TCPSegment) Parse(buffer utils.Buffer, dataGramLayerCheckSum uint32) utils.ParseResult {
	return 0
}

func (s *TCPSegment) Serialize(dataGramLayerCheckSum uint32) utils.BufferList {
	return utils.BufferList{}
}

func (s *TCPSegment) SetHeader() TCPHeader[uint32] {

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
