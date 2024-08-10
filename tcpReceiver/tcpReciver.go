package tcpReceiver

import (
	"lab/streamReassembler"
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
