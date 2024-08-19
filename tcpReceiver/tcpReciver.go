package tcpReceiver

import (
	"lab/stream"
	"lab/streamReassembler"
	"lab/tcp_helper"
	"lab/wrapping"
)

type TcpReceiverInterface interface {
}

type TcpReceiver struct {
	isn         wrapping.WrappingInt32
	setSynFlag  bool
	reassembler streamReassembler.StreamReassembler
	capacity    int
}

var _ TcpReceiverInterface = (*TcpReceiver)(nil)

func (rcv *TcpReceiver) Ackno() wrapping.WrappingInt32 {}

func (rcv *TcpReceiver) WindoSize() int {}

func (rcv *TcpReceiver) UnassembledBytes() int {}

func (rcv *TcpReceiver) SegmentReceived(seg tcp_helper.TCPSegment) {}

func (rcv *TcpReceiver) SegmentOut() stream.Stream {
	return rcv.reassembler.StreamOut()
}
